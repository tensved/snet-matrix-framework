package lib

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"regexp"
	"strings"
	"time"
)

type BotService struct {
	AIService *AIService
	Opts      BotServiceOpts
}

type BotServiceOpts struct {
	Prefix string
}

type MatrixBot interface {
	ConnectService(service *AIService, opts BotServiceOpts)
	GetService(name string) *AIService
	Run() error
}

// DefaultMatrixBot represents a Matrix bot, which can connect to multiple AI services
type DefaultMatrixBot struct {
	MauClient      *mautrix.Client // matrix client for bot
	AIServices     []BotService    // ai services
	Title          string          // title for bot
	Prefix         string          // prefix for bot, default: @
	ServicesPrefix string
	MethodsPrefix  string
	InputsPrefix   string
	startTime      time.Time
}

type MauBotOpts func(*DefaultMatrixBot)

func WithBotPrefix(prefix string) MauBotOpts {
	return func(b *DefaultMatrixBot) {
		b.Prefix = prefix
	}
}

func WithServicesPrefix(prefix string) MauBotOpts {
	return func(b *DefaultMatrixBot) {
		b.ServicesPrefix = prefix
	}
}

func WithMethodsPrefix(prefix string) MauBotOpts {
	return func(b *DefaultMatrixBot) {
		b.MethodsPrefix = prefix
	}
}

func WithInputsPrefix(prefix string) MauBotOpts {
	return func(b *DefaultMatrixBot) {
		b.InputsPrefix = prefix
	}
}

// NewDefaultMatrixBot creates a new MauBot with the specified title, client, and options
func NewDefaultMatrixBot(title string, client *mautrix.Client, opts ...MauBotOpts) *DefaultMatrixBot {
	bot := &DefaultMatrixBot{
		MauClient:      client,
		AIServices:     make([]BotService, 0),
		Title:          title,
		Prefix:         "@",
		ServicesPrefix: "-",
		MethodsPrefix:  "-",
		InputsPrefix:   "-",
		startTime:      time.Now()}

	for _, opt := range opts {
		opt(bot)
	}
	return bot
}

// ConnectService connects a service to the bot
func (b *DefaultMatrixBot) ConnectService(service *AIService, opts BotServiceOpts) {
	b.AIServices = append(b.AIServices, BotService{AIService: service, Opts: opts})
}

// Run starts the bot
func (b *DefaultMatrixBot) Run() error {

	syncer := b.MauClient.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, evt *event.Event) {

		if !b.checkEventTime(evt) || !b.checkTagMe(evt) {
			return
		}

		b.eventHandler(evt)
	})
	go func() {
		log.Info().Msgf("Starting bot %s", b.Title)
		for {
			err := b.MauClient.Sync()
			if err != nil {
				fmt.Println(err)
				time.Sleep(5 * time.Second) // Wait before retrying
			}
		}
	}()
	return nil
}

// checkEventTime checks the event time
func (b *DefaultMatrixBot) checkEventTime(evt *event.Event) bool {
	return b.startTime.Before(time.Unix(0, evt.Timestamp*int64(time.Millisecond)))
}

// checkTagMe checks if the event is tagged for the bot
func (b *DefaultMatrixBot) checkTagMe(evt *event.Event) bool {
	prefix := b.Prefix + b.Title
	return strings.HasPrefix(evt.Content.AsMessage().Body, prefix)
}

// eventHandler handles events for the bot
func (b *DefaultMatrixBot) eventHandler(evt *event.Event) {
	reMsg := regexp.MustCompile(fmt.Sprintf(`%s(?P<bot>\w+)\s+%smodel:(?P<model>\w+)\s+%smethod:(?P<method>\w+)\s*(?P<inputs>.*)`, b.Prefix, b.ServicesPrefix, b.MethodsPrefix)) // regexp for full message
	reInput := regexp.MustCompile(fmt.Sprintf(`%s(\w+):"((?:\\"|[^"])*)"`, b.InputsPrefix))                                                                                      // regexp for inputs

	msg := evt.Content.AsMessage().Body

	match := reMsg.FindStringSubmatch(msg)

	if len(match) == 0 {
		log.Debug().Msg("No match found")
		return
	}

	groups := make(map[string]string)
	for i, name := range reMsg.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = match[i]
		}
	}

	inputs, ok := groups["inputs"]
	if !ok {
		log.Debug().Msg("No inputs found")
		return
	}

	inputMatches := reInput.FindAllStringSubmatch(inputs, -1)

	params := make(map[string]interface{})

	for _, mt := range inputMatches {
		params[mt[1]] = mt[2]
	}

	svc := b.GetService(groups["model"])
	if svc == nil {
		log.Error().Msgf("Service %s not found", groups["model"])
		return
	}

	methodName, ok := groups["method"]
	if !ok {
		log.Error().Msgf("Method %s not found", groups["method"])
		return
	}

	result, err := svc.CallMethod(methodName, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call method")
		return
	}

	err = b.SendMessage(evt.RoomID, fmt.Sprintf("%v", result))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
	}
}

// GetService gets a specific AI service by name
func (b *DefaultMatrixBot) GetService(name string) *AIService {
	for _, service := range b.AIServices {
		if service.AIService.Name == name {
			return service.AIService
		}
	}
	return nil
}

// SendMessage sends a message to a specified room
func (b *DefaultMatrixBot) SendMessage(roomID id.RoomID, text string) error {

	_, err := b.MauClient.SendMessageEvent(context.Background(), roomID, event.EventMessage, event.MessageEventContent{
		MsgType:       event.MsgText,
		Body:          text,
		Format:        event.FormatHTML,
		FormattedBody: fmt.Sprintf("<p>%s</p>", text),
	})
	return err

}
