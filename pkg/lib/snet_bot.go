package lib

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"matrix-ai-framework/internal/config"
	"matrix-ai-framework/internal/matrix"
	"matrix-ai-framework/pkg/blockchain"
	"matrix-ai-framework/pkg/db"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"strings"
	"time"
)

var _ MatrixBot = (*SNETBot)(nil)

// CallState represents the state of a user interacting with the bot for sequential input filling
type CallState struct {
	ServiceName    string         // name of the called service
	Method         *AIMethod      // method that the user wants to call
	CurrentInputID int            // current input
	LastEventID    id.EventID     // last message from the bot to which a reply is expected
	FilledInputs   map[string]any // inputs, that have been filled
}

// NewCallState creates a new CallState instance
func NewCallState(serviceName string, method *AIMethod, eventID id.EventID) *CallState {
	return &CallState{
		ServiceName:    serviceName,
		Method:         method,
		CurrentInputID: 0,
		LastEventID:    eventID,
		FilledInputs:   make(map[string]any),
	}
}

// SNETBot represents a Matrix bot for SNET Services
// Implements IMauBot
type SNETBot struct {
	Client        matrix.Service
	DB            db.Service
	AIServices    []BotService
	Eth           blockchain.Ethereum
	CallStates    map[string]*CallState       // key: "{roomId} {userId}"
	PaymentStates map[string]*db.PaymentState // key: "{roomId} {userId}"
}

func NewSNETBot(client matrix.Service, db db.Service, eth blockchain.Ethereum) *SNETBot {
	return &SNETBot{
		Client:     client,
		DB:         db,
		Eth:        eth,
		CallStates: make(map[string]*CallState),
		AIServices: make([]BotService, 0),
	}
}

// ConnectService connects a service to the bot
func (bot *SNETBot) ConnectService(service *AIService, opts BotServiceOpts) {
	bot.AIServices = append(bot.AIServices, BotService{AIService: service, Opts: opts})
}

// Run starts the bot and listens for events
func (bot *SNETBot) Run() error {

	go bot.refreshToken()

	events := make(chan *event.Event)

	err := bot.Client.StartListening(events)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start matrix event listener")
		return err
	}

	for {
		select {
		case evt := <-events:
			bot.handlerEvent(evt)
		}
	}

}

// refreshToken refreshes the bot's token periodically in matrix
func (bot *SNETBot) refreshToken() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Call the Auth function every time the ticker ticks.
			bot.Client.Auth()
		}
	}
}

func (bot *SNETBot) handlerEvent(event *event.Event) {
	key := fmt.Sprintf("%s %s", event.RoomID, event.Sender)
	_, callStateFound := bot.CallStates[key]
	if !callStateFound {
		names, err := bot.parseCommand(event.Content.AsMessage().Body, event.RoomID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse command")
			return
		}

		service := bot.GetService(names.ServiceName)
		if service == nil {
			log.Error().Msg("Service not found")
			return
		}
		method := service.GetMethod(names.MethodName)
		if method == nil {
			log.Error().Msg("Method not found")
			return
		}
		inputs := method.Inputs
		if len(inputs) == 0 {
			log.Error().Msg("No inputs")
			return
		}

		paymentStateKey := fmt.Sprintf("%s %s %s %s", event.RoomID, event.Sender, names.ServiceName, names.MethodName)

		paymentStateUUID := uuid.New()
		recipientAddress := config.Blockchain.AdminPublicAddress

		nameParts := strings.Split(names.ServiceName, "/")
		s, err := bot.DB.GetSnetService(nameParts[0])
		if err != nil {
			log.Error().Err(err).Msg("Failed to get snet service")
		}
		tokenAddress, err := bot.Eth.MPE.Token(&bind.CallOpts{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to get token address")
		}
		paymentStateURL := fmt.Sprintf("ethereum:%s/transfer?address=%s&uint256=%s",
			tokenAddress, recipientAddress, s.Price)

		paymentState := &db.PaymentState{
			ID:           paymentStateUUID,
			URL:          paymentStateURL,
			Key:          paymentStateKey,
			TokenAddress: tokenAddress.Hex(),
			ToAddress:    recipientAddress,
			Amount:       s.Price,
		}
		_, err = bot.DB.CreatePaymentState(paymentState)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create payment state")
			return
		}
		paymentGatewayURL := fmt.Sprintf("%s?id=%s", fmt.Sprintf("https://%s", config.App.Domain), paymentStateUUID.String())
		text := fmt.Sprintf("To continue, pay for the service at the link:<br>%s", paymentGatewayURL)

		_, err = bot.Client.SendMessage(event.RoomID, text)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}

		go bot.waitForPayment(event, paymentStateUUID, names)
	} else {
		bot.processCallState(event, ParsedNames{})
	}
}

func (bot *SNETBot) waitForPayment(event *event.Event, paymentID uuid.UUID, names ParsedNames) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			paymentState, err := bot.DB.GetPaymentState(paymentID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get payment state")
				continue
			}

			txHash := paymentState.TxHash
			if txHash == nil {
				log.Error().Msg("TxHash not found in payment state")
				continue
			}

			receipt, err := bot.Eth.Client.TransactionReceipt(context.Background(), common.HexToHash(*txHash))
			if err != nil {
				log.Error().Err(err).Msg("Failed to get transaction receipt")
				continue
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				paymentState.Status = "paid"
				err = bot.DB.PatchUpdatePaymentState(paymentState)
				if err != nil {
					log.Error().Err(err).Msg("Failed to update payment state status")
					continue
				}

				text := "Payment received. You can now proceed with the service."
				_, err = bot.Client.SendMessage(event.RoomID, text)
				if err != nil {
					log.Error().Err(err).Msg("Failed to send message")
					continue
				}

				bot.processCallState(event, names)
				return
			}
		case <-time.After(10 * time.Minute):
			log.Info().Msg("Timeout waiting for payment")
			return
		}
	}
}

func (bot *SNETBot) processCallState(event *event.Event, names ParsedNames) {
	key := fmt.Sprintf("%s %s", event.RoomID, event.Sender)
	callState, callStateFound := bot.CallStates[key]
	if !callStateFound {
		service := bot.GetService(names.ServiceName)
		if service == nil {
			log.Error().Msg("Service not found")
			return
		}

		method := service.GetMethod(names.MethodName)
		if method == nil {
			log.Error().Msg("Method not found")
			return
		}

		inputs := method.Inputs
		if len(inputs) == 0 {
			log.Error().Msg("No inputs")
			return
		}

		text := bot.getInputDescription(inputs[0])

		resp, err := bot.Client.SendMessage(event.RoomID, text)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
			return
		}

		callState = NewCallState(names.ServiceName, method, resp.EventID)
		bot.CallStates[key] = callState

		return
	}

	repliedEvt, err := bot.Client.GetRepliedEvent(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get replied event")
		return
	}

	if callState.LastEventID != repliedEvt.ID {
		return
	}

	currentInput := callState.Method.Inputs[callState.CurrentInputID]

	msg := event.Content.AsMessage()

	_, replyText, err := matrix.ExtractTexts(msg.FormattedBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract text")
		return
	}

	log.Debug().Msgf("Got new message: %v", replyText)

	callState.FilledInputs[currentInput.Name] = replyText

	callState.CurrentInputID++

	if callState.CurrentInputID >= len(callState.Method.Inputs) {
		log.Debug().Msgf("Full data filled: %v", callState.FilledInputs)

		ctx := NewMContext(WithParams(callState.FilledInputs))

		go callState.Method.Handler.Call(ctx)

		_, err = bot.Client.SendMessage(event.RoomID, "Wait for answer...")
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
			return
		}

		delete(bot.CallStates, key)

		select {
		case res := <-ctx.Result:
			log.Debug().Msgf("Got res: %v", res)
			_, err = bot.Client.SendMessage(event.RoomID, fmt.Sprintf("<code>%s</code>", res))
			if err != nil {
				log.Error().Err(err).Msg("Failed to send message")
			}

		case <-time.After(2 * time.Minute):
			log.Debug().Msg("Timeout!")
		}

		return
	}

	text := bot.getInputDescription(callState.Method.Inputs[callState.CurrentInputID])
	resp, err := bot.Client.SendMessage(event.RoomID, text)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return
	}

	callState.LastEventID = resp.EventID
}

// ParsedNames contains parsed information from a command
type ParsedNames struct {
	BotName     string // bot name. If private chat, it will be empty
	ServiceName string
	MethodName  string
}

// parseCommand parses the command message and extracts relevant information
func (bot *SNETBot) parseCommand(message string, roomID id.RoomID) (ParsedNames, error) {
	isPrivate, err := bot.Client.IsPrivateRoom(roomID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get room info")
		return ParsedNames{}, err
	}

	trimmedMessage := strings.TrimSpace(message)

	names := strings.Split(trimmedMessage, " ")

	if !isPrivate {
		if len(names) != 3 {
			return ParsedNames{}, fmt.Errorf("incorrect params number. Want 3, got %d", len(names))
		}

		return ParsedNames{
			BotName:     names[0],
			ServiceName: names[1],
			MethodName:  names[2],
		}, nil
	}

	if len(names) != 2 {
		return ParsedNames{}, fmt.Errorf("incorrect params number. Want 2, got %d", len(names))
	}

	return ParsedNames{
		ServiceName: names[0],
		MethodName:  names[1],
	}, nil
}

// getInputDescription returns a description for the given input
func (bot *SNETBot) getInputDescription(input MInput) string {
	return fmt.Sprintf("Input param %s", input.Name)
}

// GetService retrieves a specific AI service by name
func (bot *SNETBot) GetService(name string) *AIService {
	for _, service := range bot.AIServices {
		if service.AIService.Name == name {
			return service.AIService
		}
	}
	return nil
}
