package snet

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix"
	"github.com/tensved/bobrix/contracts"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/internal/matrix"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"google.golang.org/grpc/health/grpc_health_v1"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"strings"
	"time"
)

// ParsedNames contains parsed information from a user's message in Matrix room.
type ParsedNames struct {
	Bot        string // bot name. If private chat, it will be empty
	SnetID     string
	Descriptor string
	Service    string
	Method     string
}

// CallState represents the state of a user interacting with the bot for sequential input filling.
type CallState struct {
	ServiceName    string            // name of the called service
	Method         *contracts.Method // method that the user wants to call
	CurrentInputID int               // current input
	LastEventID    id.EventID        // last message from the bot to which a reply is expected
	FilledInputs   map[string]any    // inputs, that have been filled
}

// NewCallState creates a new CallState instance.
func NewCallState(serviceName string, method *contracts.Method, eventID id.EventID) *CallState {
	return &CallState{
		ServiceName:    serviceName,
		Method:         method,
		CurrentInputID: 0,
		LastEventID:    eventID,
		FilledInputs:   make(map[string]any),
	}
}

func Parser(mx matrix.Service, eth blockchain.Ethereum, database db.Service, callStates map[string]*CallState, bobr *bobrix.Bobrix, grpc *grpcmanager.GRPCClientManager) func(evt *event.Event) *bobrix.AIRequest {
	return func(evt *event.Event) *bobrix.AIRequest {
		key := fmt.Sprintf("%s %s", evt.RoomID, evt.Sender)

		callState, callStateFound := callStates[key]
		if !callStateFound {
			names, err := parseCommand(evt.Content.AsMessage().Body, evt.RoomID, mx)
			if err != nil {
				log.Error().Err(err)
				return nil
			}

			snetService, err := database.GetSnetService(names.SnetID)
			if err != nil {
				log.Error().Err(err)
				_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
				if err != nil {
					log.Error().Err(err)
					return nil
				}
				return nil
			}
			if snetService == nil || snetService.URL == "" {
				log.Error().Msgf("Service with ID %s not found", names.SnetID)
				_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
				if err != nil {
					log.Error().Err(err)
					return nil
				}
				return nil
			}

			client, err := grpc.GetClient(snetService.URL)
			if err != nil {
				log.Error().Err(err)
				_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
				if err != nil {
					log.Error().Err(err)
					return nil
				}
				return nil
			}

			healthClient := grpc_health_v1.NewHealthClient(client.Conn)
			hReq := grpc_health_v1.HealthCheckRequest{}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			hResp, err := healthClient.Check(ctx, &hReq)
			if err != nil {
				log.Error().Err(err).Msgf("failed to get health status: %+v", hResp)
				_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
				if err != nil {
					log.Error().Err(err)
					return nil
				}
				return nil
			}

			if hResp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
				log.Debug().Msgf("service is offline: %+v", hResp.GetStatus())
				return nil
			}
			log.Debug().Msgf("service is online: %+v", hResp.GetStatus())

			paymentStateKey := fmt.Sprintf("%s %s %s", evt.RoomID, evt.Sender, evt.Content.AsMessage().Body)

			service, found := bobr.GetService(names.SnetID)
			if !found {
				log.Error().Err(errors.New("service not found"))
				_, err = mx.SendMessage(evt.RoomID, "Service unavailable.")
				if err != nil {
					log.Error().Err(err)
					return nil
				}
				return nil
			}

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

			paymentStateUUID := uuid.New()
			recipientAddress := config.Blockchain.AdminPublicAddress

			tokenAddress, err := eth.MPE.Token(&bind.CallOpts{})
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			paymentStateURL := fmt.Sprintf("ethereum:%s/transfer?address=%s&uint256=%v",
				tokenAddress, recipientAddress, snetService.Price)

			paymentState := &db.PaymentState{
				ID:           paymentStateUUID,
				URL:          paymentStateURL,
				Key:          paymentStateKey,
				TokenAddress: tokenAddress.Hex(),
				ToAddress:    recipientAddress,
				Amount:       snetService.Price,
			}
			_, err = database.CreatePaymentState(paymentState)
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			paymentGatewayURL := fmt.Sprintf("%s?id=%s", fmt.Sprintf("http://%s", config.App.Domain), paymentStateUUID.String())
			text := fmt.Sprintf("To continue, pay for the service within %d min:<br>%s", config.App.PaymentTimeout, paymentGatewayURL)

			_, err = mx.SendMessage(evt.RoomID, text)
			if err != nil {
				log.Error().Err(err)
				return nil
			}

			err = waitForPayment(evt, paymentStateUUID, mx, eth, database)
			if err != nil {
				log.Error().Err(err)
				return nil
			}

			fieldDescription := method.Inputs[0].Description
			resp, err := mx.SendMessage(evt.RoomID, fieldDescription)
			if err != nil {
				log.Error().Err(err)
				return nil
			}
			callStates[key] = NewCallState(names.SnetID, method, resp.EventID)
		} else {
			err := processCallState(evt, mx, callState)
			if err != nil {
				log.Error().Err(err)
				return nil
			}

			if callStates[key].CurrentInputID >= len(callStates[key].Method.Inputs) {
				log.Debug().Msgf("full data filled: %v", callStates[key].FilledInputs)

				req := &bobrix.AIRequest{
					ServiceName: callStates[key].ServiceName,
					MethodName:  callStates[key].Method.Name,
					InputParams: callStates[key].FilledInputs,
				}

				delete(callStates, key)

				return req
			}
		}

		return nil
	}
}

// parseCommand parses the command message and extracts relevant information.
func parseCommand(msg string, roomID id.RoomID, mx matrix.Service) (ParsedNames, error) {
	isPrivate, err := mx.IsPrivateRoom(roomID)
	if err != nil {
		return ParsedNames{}, err
	}

	trimmed := strings.TrimSpace(msg)

	names := strings.Split(trimmed, " ")

	privNamesNumber := 4
	pubNamesNumber := 5
	if !isPrivate {
		if len(names) != pubNamesNumber {
			return ParsedNames{}, fmt.Errorf("incorrect params number. Want 5, got %d", len(names))
		}

		return ParsedNames{
			Bot:        names[0],
			SnetID:     names[1],
			Descriptor: names[2],
			Service:    names[3],
			Method:     names[4],
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
	}, nil
}

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

	text := callState.Method.Inputs[callState.CurrentInputID].Description
	resp, err := mx.SendMessage(event.RoomID, text)
	if err != nil {
		return err
	}

	callState.LastEventID = resp.EventID

	return nil
}
