package lib

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matrix-ai-framework/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"strings"
	"time"
)

var _ MatrixBot = (*SNETBot)(nil)

// UserState represents the state of a user interacting with the bot for sequential input filling
type UserState struct {
	ServiceName    string         // name of the called service
	Method         *AIMethod      // method that the user wants to call
	CurrentInputID int            // current input
	LastEventID    id.EventID     // last message from the bot to which a reply is expected
	FilledInputs   map[string]any // inputs, that have been filled
}

// NewUserState creates a new UserState instance
func NewUserState(serviceName string, method *AIMethod, eventID id.EventID) *UserState {
	return &UserState{
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
	Client     matrix.Service
	AIServices []BotService
	States     map[string]*UserState // key: "{roomId} {userId}"
}

func NewSNETBot(client matrix.Service) *SNETBot {
	return &SNETBot{
		Client:     client,
		States:     make(map[string]*UserState),
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

// handlerEvent handles incoming events and user interactions
func (bot *SNETBot) handlerEvent(event *event.Event) {
	key := fmt.Sprintf("%s %s", event.RoomID, event.Sender)

	state, found := bot.States[key]
	if !found {

		names, err := bot.parseCommand(event.Content.AsMessage().Body, event.RoomID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse command: " + err.Error())
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

		text := bot.getInputDescription(inputs[0])

		resp, err := bot.Client.SendMessage(event.RoomID, text)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message: " + err.Error())
			return
		}

		state = NewUserState(names.ServiceName, method, resp.EventID)
		bot.States[key] = state

		return
	}

	repliedEvt, err := bot.Client.GetRepliedEvent(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get replied event")
		return
	}

	if state.LastEventID != repliedEvt.ID {
		return
	}

	currentInput := state.Method.Inputs[state.CurrentInputID]

	msg := event.Content.AsMessage()

	_, replyText, err := matrix.ExtractTexts(msg.FormattedBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract text")
		return
	}

	log.Info().Msgf("GET NEW MSG: %v", replyText)

	state.FilledInputs[currentInput.Name] = replyText

	state.CurrentInputID++

	if state.CurrentInputID >= len(state.Method.Inputs) {

		log.Info().Msg("Full data filled: " + fmt.Sprint(state.FilledInputs))

		ctx := NewMContext(WithParams(state.FilledInputs))

		go state.Method.Handler.Call(ctx)

		_, err := bot.Client.SendMessage(event.RoomID, "Wait for answer...")
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message: " + err.Error())
			return
		}

		delete(bot.States, key)

		select {
		case res := <-ctx.Result:
			log.Debug().Msgf("Got res: %v", res)
			_, err := bot.Client.SendMessage(event.RoomID, fmt.Sprintf("<code>%s</code>", res))
			if err != nil {
				log.Error().Err(err).Msg("Failed to send message: " + err.Error())
			}

		case <-time.After(5 * time.Second):
			log.Info().Msg("Timeout!")
		}

		return
	}

	text := bot.getInputDescription(state.Method.Inputs[state.CurrentInputID])
	resp, err := bot.Client.SendMessage(event.RoomID, text)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message: " + err.Error())
		return
	}

	state.LastEventID = resp.EventID
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

	names := strings.Split(message, " ")

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
