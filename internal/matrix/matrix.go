package matrix

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/internal/syncer"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"golang.org/x/net/html"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Service defines the methods for interacting with the Matrix Synapse server.
type Service interface {
	Register(username, password string) (err error)
	Login(username, password string) (err error)
	Auth()
	SendMessage(roomID id.RoomID, text string) (*mautrix.RespSendEvent, error)
	GetRepliedEvent(evt *event.Event) (*event.Event, error)
	IsPrivateRoom(roomID id.RoomID) (bool, error)
}

// service is an implementation of the Service interface.
type service struct {
	Client      *mautrix.Client                // Client is the Matrix client used to communicate with the Matrix server.
	Context     context.Context                // Context is used for managing request lifecycles.
	Syncer      *mautrix.DefaultSyncer         // Syncer is used to synchronize events with the Matrix server.
	startTime   time.Time                      // startTime holds the service start time.
	db          db.Service                     // db provides access to the database layer.
	snetSyncer  syncer.SnetSyncer              // snetSyncer is responsible for network synchronization.
	grpcManager *grpcmanager.GRPCClientManager // grpcManager manages gRPC client connections.
	eth         blockchain.Ethereum            // eth provides access to the Ethereum blockchain.
}

// New creates a new instance of the service and initializes the Matrix client.
func New(db db.Service, snetSyncer syncer.SnetSyncer, grpcManager *grpcmanager.GRPCClientManager, eth blockchain.Ethereum) Service {
	logger := log.With().
		Str("homeserver", config.Matrix.HomeserverURL).
		Str("username", config.Matrix.Username).
		Logger()

	client, err := mautrix.NewClient(config.Matrix.HomeserverURL, "", "")
	if err != nil {
		logger.Error().Err(err).Msg("failed to create Matrix client")
		return nil
	}

	sync, ok := client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		logger.Error().Msg("failed to assert client.Syncer to *mautrix.DefaultSyncer")
		return nil
	}

	m := &service{client, context.Background(), sync, time.Now(), db, snetSyncer, grpcManager, eth}
	logger.Info().Msg("Matrix client created successfully")
	return m
}

// Register registers a new user with the given username and password on the Matrix server.
func (s *service) Register(username, password string) error {
	logger := log.With().Str("username", username).Logger()

	resp, err := s.Client.RegisterDummy(s.Context, &mautrix.ReqRegister{
		Username:     username,
		Password:     password,
		InhibitLogin: false,
		Auth:         nil,
		Type:         "m.login.password",
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to register Matrix user")
		return err
	}

	s.Client.UserID = resp.UserID
	s.Client.AccessToken = resp.AccessToken

	logger.Info().
		Str("user_id", string(resp.UserID)).
		Msg("Matrix user registered successfully")
	return nil
}

// Login logs in a user with the given username and password on the Matrix server.
func (s *service) Login(username, password string) error {
	logger := log.With().Str("username", username).Logger()

	resp, err := s.Client.Login(s.Context, &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: username,
		},
		Password: password,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to login Matrix user")
		return err
	}

	s.Client.UserID = resp.UserID
	s.Client.AccessToken = resp.AccessToken

	logger.Info().
		Str("user_id", string(resp.UserID)).
		Msg("Matrix user logged in successfully")
	return nil
}

// GetUserProfile retrieves the profile of a user by their user ID.
func (s *service) GetUserProfile(userID string) error {
	logger := log.With().Str("user_id", userID).Logger()

	matrixUserID := id.UserID(userID)
	_, err := s.Client.GetProfile(s.Context, matrixUserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get Matrix user profile")
		return err
	}

	logger.Debug().Msg("Matrix user profile retrieved successfully")
	return nil
}

// Auth performs authentication by either logging in or registering the user depending on their existing profile.
func (s *service) Auth() {
	logger := log.With().
		Str("username", config.Matrix.Username).
		Str("homeserver", config.Matrix.HomeserverURL).
		Logger()

	userID := fmt.Sprintf("@%s:%s", config.Matrix.Username, config.Matrix.Servername)
	err := s.GetUserProfile(userID)
	if err != nil {
		logger.Info().
			Str("user_id", userID).
			Err(err).
			Msg("user profile not found, attempting registration")
		err = s.Register(config.Matrix.Username, config.Matrix.Password)
		if err != nil {
			logger.Error().
				Err(err).
				Str("username", config.Matrix.Username).
				Msg("failed to register user during auth")
			return
		}
		logger.Info().
			Str("username", config.Matrix.Username).
			Msg("user registered successfully during auth")
	} else {
		logger.Info().
			Str("user_id", userID).
			Msg("user profile found, attempting login")
		err = s.Login(config.Matrix.Username, config.Matrix.Password)
		if err != nil {
			logger.Error().
				Err(err).
				Str("username", config.Matrix.Username).
				Msg("failed to login user during auth")
			return
		}
		logger.Info().
			Str("username", config.Matrix.Username).
			Msg("user logged in successfully during auth")
	}

	logger.Info().Msg("Matrix authentication completed successfully")
}

// isReply checks if the event is a reply to another event.
func isReply(evt *event.Event) bool {
	return evt.Content.AsMessage().RelatesTo != nil && evt.Content.AsMessage().RelatesTo.InReplyTo != nil
}

// IsPrivateRoom checks if the given room ID corresponds to a private room with exactly two members.
func (s *service) IsPrivateRoom(roomID id.RoomID) (bool, error) {
	members, err := s.Client.JoinedMembers(s.Context, roomID)
	if err != nil {
		return false, err
	}
	return len(members.Joined) == 2, nil
}

// GetRepliedEvent retrieves the event to which the given event is a reply.
func (s *service) GetRepliedEvent(evt *event.Event) (*event.Event, error) {
	if !isReply(evt) {
		return nil, errors.New("not a reply")
	}
	return s.Client.GetEvent(s.Context, evt.RoomID, evt.Content.AsMessage().RelatesTo.InReplyTo.EventID)
}

// ExtractTexts extracts the original and reply texts from a formatted body of an event.
func ExtractTexts(formattedBody string) (originalText, replyText string, err error) {
	doc, err := html.Parse(strings.NewReader(formattedBody))
	if err != nil {
		return "", "", err
	}

	var extract func(*html.Node, bool)
	extract = func(n *html.Node, insideMXReply bool) {
		if n.Type == html.ElementNode && n.Data == "mx-reply" {
			insideMXReply = true
		} else if n.Type == html.ElementNode && n.Data == "br" && insideMXReply {
			originalText = ""
		}

		if n.Type == html.TextNode && insideMXReply {
			originalText += n.Data
		} else if n.Type == html.TextNode && !insideMXReply {
			replyText += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c, insideMXReply)
		}
	}

	extract(doc, false)
	return strings.TrimSpace(originalText), strings.TrimSpace(replyText), nil
}

// SendMessage sends a text message to the specified room on the Matrix server.
func (s *service) SendMessage(roomID id.RoomID, text string) (*mautrix.RespSendEvent, error) {
	logger := log.With().
		Str("room_id", string(roomID)).
		Str("message_length", fmt.Sprintf("%d", len(text))).
		Logger()

	content := event.MessageEventContent{
		MsgType:       event.MsgText,
		Body:          text,
		Format:        event.FormatHTML,
		FormattedBody: fmt.Sprintf("<p>%s</p>", text),
	}

	resp, err := s.Client.SendMessageEvent(s.Context, roomID, event.EventMessage, content)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to send Matrix message")
		return nil, err
	}

	logger.Debug().
		Str("event_id", string(resp.EventID)).
		Msg("Matrix message sent successfully")

	return resp, nil
}
