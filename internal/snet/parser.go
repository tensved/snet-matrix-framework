package snet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix"
	"github.com/tensved/bobrix/contracts"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/internal/matrix"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	ipfsutils "github.com/tensved/snet-matrix-framework/pkg/ipfs"
	"google.golang.org/grpc/health/grpc_health_v1"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// ParsedNames contains parsed information from a user's message in Matrix room.
type ParsedNames struct {
	Bot        string                 // bot name. If private chat, it will be empty
	SnetID     string                 // service ID
	Descriptor string                 // service descriptor
	Service    string                 // service name
	Method     string                 // method name
	Params     map[string]interface{} // method parameters in JSON format
}

// Deprecated: CallState represents the state of a user interacting with the bot for sequential input filling.
type CallState struct {
	ServiceName    string            // name of the called service
	Method         *contracts.Method // method that the user wants to call
	CurrentInputID int               // current input
	LastEventID    id.EventID        // last message from the bot to which a reply is expected
	FilledInputs   map[string]any    // inputs, that have been filled
}

// Deprecated: NewCallState creates a new CallState instance.
func NewCallState(serviceName string, method *contracts.Method, eventID id.EventID) *CallState {
	return &CallState{
		ServiceName:    serviceName,
		Method:         method,
		CurrentInputID: 0,
		LastEventID:    eventID,
		FilledInputs:   make(map[string]any),
	}
}

func Parser(mx matrix.Service, eth blockchain.Ethereum, database db.Service, callStates map[string]*CallState, bobr *bobrix.Bobrix, grpc *grpcmanager.GRPCClientManager) func(evt *event.Event) *bobrix.ServiceRequest {
	return func(evt *event.Event) *bobrix.ServiceRequest {
		// Skip if message starts with ! (bot commands)
		if strings.HasPrefix(strings.TrimSpace(evt.Content.AsMessage().Body), "!") {
			return nil
		}

		names, err := parseCommand(evt.Content.AsMessage().Body, evt.RoomID, mx)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse command")
			return nil
		}

		log.Info().Str("snet_id", names.SnetID).Str("descriptor", names.Descriptor).Str("service", names.Service).Str("method", names.Method).Interface("params", names.Params).Msg("parsed command")

		snetService, err := database.GetSnetService(names.SnetID)
		if err != nil {
			log.Error().Err(err).Str("snet_id", names.SnetID).Msg("failed to get service from database")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		if snetService == nil || snetService.URL == "" {
			log.Error().Str("snet_id", names.SnetID).Msg("service not found in database or URL is empty")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}

		log.Info().Str("snet_id", names.SnetID).Str("url", snetService.URL).Msg("found service in database")

		log.Info().Str("url", snetService.URL).Msg("attempting to get gRPC client")
		client, err := grpc.GetClient(snetService.URL)
		if err != nil {
			log.Error().Err(err).Str("url", snetService.URL).Msg("failed to get gRPC client")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Str("url", snetService.URL).Msg("successfully got gRPC client")

		log.Info().Str("url", snetService.URL).Msg("checking service health")
		healthClient := grpc_health_v1.NewHealthClient(client.Conn)
		hReq := grpc_health_v1.HealthCheckRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hResp, err := healthClient.Check(ctx, &hReq)
		if err != nil {
			log.Error().Err(err).Str("url", snetService.URL).Msg("failed to get health status")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}

		log.Info().Str("url", snetService.URL).Str("status", hResp.GetStatus().String()).Msg("health check response")
		if hResp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
			log.Error().Str("url", snetService.URL).Str("status", hResp.GetStatus().String()).Msg("service is offline")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Str("url", snetService.URL).Msg("service is online")

		log.Info().Str("snet_id", names.SnetID).Msg("checking service in bobrix")
		service, found := bobr.GetService(names.SnetID)
		if !found {
			log.Error().Str("snet_id", names.SnetID).Msg("service not found in bobrix")
			_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Str("snet_id", names.SnetID).Msg("found service in bobrix")

		method := service.Service.Methods[names.Method]
		if method == nil {
			log.Error().Err(errors.New("method not found"))
			_, err = mx.SendMessage(evt.RoomID, "Method unavailable.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}

		log.Info().
			Str("snet_id", names.SnetID).
			Str("method", names.Method).
			Msg("using new payment system")

		privateKey, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse private key")
			_, err = mx.SendMessage(evt.RoomID, "Internal error.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Msg("private key parsed successfully")

		protoFiles, err := getProtoFilesForService(snetService)
		if err != nil {
			log.Error().Err(err).Msg("failed to get proto files")
			_, err = mx.SendMessage(evt.RoomID, "Internal error.")
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Msg("proto files obtained successfully")

		paymentManager := NewPaymentManager(eth, database, grpc, privateKey, protoFiles)
		log.Info().Msg("payment manager created successfully")

		log.Info().Msg("calling PaymentManager.ExecuteCall")
		result, err := paymentManager.ExecuteCall(context.Background(), snetService, names.Method, names.Params)
		if err != nil {
			log.Error().Err(err).Msg("failed to execute call")
			_, err = mx.SendMessage(evt.RoomID, fmt.Sprintf("Error: %v", err))
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			return nil
		}
		log.Info().Msg("PaymentManager.ExecuteCall completed successfully")

		resultStr := fmt.Sprintf("%v", result)
		_, err = mx.SendMessage(evt.RoomID, resultStr)
		if err != nil {
			log.Error().Err(err)
			return nil
		}
		log.Info().Msg("result sent to Matrix successfully")

		return nil
	}
}

// parseCommand parses the command message and extracts relevant information.
func parseCommand(msg string, roomID id.RoomID, mx matrix.Service) (ParsedNames, error) {
	logger := log.With().
		Str("message", msg).
		Str("room_id", string(roomID)).
		Logger()

	isPrivate, err := mx.IsPrivateRoom(roomID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to check if room is private")
		return ParsedNames{}, err
	}

	trimmed := strings.TrimSpace(msg)

	// Try to find JSON parameters at the end of the message
	// Format: "snet_id descriptor service method {json_params}"
	lastBraceIndex := strings.LastIndex(trimmed, "{")
	lastClosingBraceIndex := strings.LastIndex(trimmed, "}")

	var params map[string]interface{}
	var commandPart string

	if lastBraceIndex != -1 && lastClosingBraceIndex != -1 && lastClosingBraceIndex > lastBraceIndex {
		jsonStr := trimmed[lastBraceIndex : lastClosingBraceIndex+1]
		commandPart = strings.TrimSpace(trimmed[:lastBraceIndex])

		err := json.Unmarshal([]byte(jsonStr), &params)
		if err != nil {
			logger.Error().
				Err(err).
				Str("json_params", jsonStr).
				Msg("failed to parse JSON parameters")
			return ParsedNames{}, fmt.Errorf("invalid JSON parameters: %w", err)
		}
	} else {
		commandPart = trimmed
		params = make(map[string]interface{})
		logger.Debug().Msg("no JSON parameters found in command")
	}

	names := strings.Split(commandPart, " ")

	privNamesNumber := 4
	pubNamesNumber := 5
	if !isPrivate {
		if len(names) != pubNamesNumber {
			logger.Error().
				Int("expected", pubNamesNumber).
				Int("got", len(names)).
				Msg("incorrect number of parameters for public room")
			return ParsedNames{}, fmt.Errorf("incorrect params number. Want 5, got %d", len(names))
		}

		return ParsedNames{
			Bot:        names[0],
			SnetID:     names[1],
			Descriptor: names[2],
			Service:    names[3],
			Method:     names[4],
			Params:     params,
		}, nil
	}

	if len(names) != privNamesNumber {
		return ParsedNames{}, fmt.Errorf("incorrect params number. Want 4, got %d", len(names))
	}

	return ParsedNames{
		SnetID:     names[0],
		Descriptor: names[1],
		Service:    names[2],
		Method:     names[3],
		Params:     params,
	}, nil
}

// getProtoFilesForService retrieves proto files for the service from IPFS
func getProtoFilesForService(snetService *db.SnetService) (map[string]string, error) {
	ipfsClient := ipfsutils.Init()

	var protoHashes []string

	if snetService.ServiceApiSource != "" {
		protoHashes = append(protoHashes, snetService.ServiceApiSource)
	}

	if snetService.ModelIpfsHash != "" {
		protoHashes = append(protoHashes, snetService.ModelIpfsHash)
	}

	if len(protoHashes) == 0 {
		return nil, fmt.Errorf("both ModelIpfsHash and ServiceApiSource are empty")
	}

	var content []byte
	var protoErr error

	for _, protoHash := range protoHashes {
		log.Info().Str("hash", protoHash).Msg("trying to get proto files from IPFS")
		content, protoErr = ipfsClient.GetIpfsFile(protoHash)
		if protoErr == nil {
			log.Info().Str("successful_hash", protoHash).Msg("successfully got proto files from IPFS")
			break
		}
		log.Error().Err(protoErr).Str("hash", protoHash).Msg("failed to get proto files from IPFS, trying next hash")
	}

	if protoErr != nil {
		return nil, fmt.Errorf("failed to get proto files from all IPFS hashes: %w", protoErr)
	}

	protoFiles, err := ipfsutils.ReadFilesCompressed(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed proto files: %w", err)
	}

	protoFilesMap := make(map[string]string)
	for fileName, fileContent := range protoFiles {
		protoFilesMap[fileName] = string(fileContent)
	}

	log.Info().Int("files_count", len(protoFilesMap)).Msg("extracted proto files count")
	return protoFilesMap, nil
}

// Deprecated: waitForPayment waits for payment confirmation.
func waitForPayment(event *event.Event, paymentID uuid.UUID, mx matrix.Service, eth blockchain.Ethereum, database db.Service) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(time.Duration(config.App.PaymentTimeout) * time.Minute)

	for {
		select {
		case <-ticker.C:
			paymentState, err := database.GetPaymentState(paymentID)
			if err != nil {
				log.Error().Err(err)
				continue
			}

			txHash := paymentState.TxHash
			if txHash == nil {
				log.Debug().Msg("tx hash not found in payment state")
				continue
			}

			receipt, err := eth.Client.TransactionReceipt(context.Background(), common.HexToHash(*txHash))
			if err != nil {
				log.Error().Err(err)
				continue
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				paymentState.Status = "paid"
				err = database.PatchUpdatePaymentState(paymentState)
				if err != nil {
					return err
				}

				text := "Payment has been received. Now you can reply to the bot messages to complete the input fields."
				_, err = mx.SendMessage(event.RoomID, text)
				if err != nil {
					return err
				}
				return nil
			}
		case <-timeout:
			log.Debug().Msg("payment timeout reached")

			err := database.PatchUpdatePaymentState(&db.PaymentState{Status: "expired"})
			if err != nil {
				return err
			}

			_, err = mx.SendMessage(event.RoomID, "Waiting time for payment has expired. Please, try again.")
			if err != nil {
				return err
			}

			return errors.New("waiting time for payment has expired")
		}
	}
}

// Deprecated: processCallState processes the call state for sequential input filling.
func processCallState(event *event.Event, mx matrix.Service, callState *CallState) error {
	repliedEvt, err := mx.GetRepliedEvent(event)
	if err != nil {
		return err
	}

	if callState.LastEventID != repliedEvt.ID {
		return errors.New("incorrect last event id")
	}

	currentInput := callState.Method.Inputs[callState.CurrentInputID]

	msg := event.Content.AsMessage()

	_, replyText, err := matrix.ExtractTexts(msg.FormattedBody)
	if err != nil {
		return err
	}

	callState.FilledInputs[currentInput.Name] = replyText

	callState.CurrentInputID++

	if callState.CurrentInputID >= len(callState.Method.Inputs) {
		log.Debug().Msgf("full data filled: %+v", callState.FilledInputs)
		return nil
	}

	textMap := callState.Method.Inputs[callState.CurrentInputID].Description
	nextText := textMap["en"]
	if nextText == "" {
		for _, v := range textMap {
			nextText = v
			break
		}
	}
	resp, err := mx.SendMessage(event.RoomID, nextText)
	if err != nil {
		return err
	}

	callState.LastEventID = resp.EventID

	return nil
}

// Deprecated: getPrivateKey retrieves the private key from configuration
func getPrivateKey() *ecdsa.PrivateKey {
	privateKeyECDSA, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse private key")
		return nil
	}
	return privateKeyECDSA
}
