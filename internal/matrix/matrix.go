package matrix

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"matrix-ai-framework/internal/config"
	"matrix-ai-framework/internal/grpc_manager"
	"matrix-ai-framework/internal/snet_syncer"
	"matrix-ai-framework/pkg/blockchain"
	"matrix-ai-framework/pkg/db"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	Register(username, password string) (err error)
	Login(username, password string) (err error)
	Auth()
	StartListening(ch chan *event.Event) (err error)
	SendMessage(roomID id.RoomID, text string) (*mautrix.RespSendEvent, error)
	SendMessageWithMedia(roomID, mxcURI, message string) error
	IsPrivateRoom(roomID id.RoomID) (bool, error)
	GetRepliedEvent(evt *event.Event) (*event.Event, error)
}

type service struct {
	Client      *mautrix.Client
	Context     context.Context
	Syncer      *mautrix.DefaultSyncer
	startTime   time.Time
	db          db.Service
	snetSyncer  snet_syncer.SnetSyncer
	grpcManager *grpc_manager.GRPCClientManager
	eth         blockchain.Ethereum
}

func New(db db.Service, snetSyncer snet_syncer.SnetSyncer, grpcManager *grpc_manager.GRPCClientManager, eth blockchain.Ethereum) Service {
	client, err := mautrix.NewClient(config.Matrix.HomeserverURL, "", "")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Matrix client")
	}
	ctx := context.Background()
	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	m := &service{client, ctx, syncer, time.Now(), db, snetSyncer, grpcManager, eth}
	log.Info().Msg("Matrix connect established")
	return m
}

func (s *service) Register(username, password string) (err error) {
	resp, err := s.Client.RegisterDummy(s.Context, &mautrix.ReqRegister{
		Username:     username,
		Password:     password,
		InhibitLogin: false,
		Auth:         nil,
		Type:         "m.login.password",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to register")
	}
	s.Client.UserID = resp.UserID
	s.Client.AccessToken = resp.AccessToken
	return
}

func (s *service) Login(username, password string) (err error) {
	resp, err := s.Client.Login(s.Context, &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: username,
		},
		Password: password,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to login")
	}
	s.Client.UserID = resp.UserID
	s.Client.AccessToken = resp.AccessToken
	return
}

func (s *service) GetUserProfile(userID string) (err error) {
	matrixUserID := id.UserID(userID)
	_, err = s.Client.GetProfile(s.Context, matrixUserID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get user profile")
	}
	return
}

func (s *service) Auth() {
	userID := fmt.Sprintf("@%s:%s", config.Matrix.Username, config.Matrix.Servername)
	err := s.GetUserProfile(userID)
	if err != nil {
		err = s.Register(config.Matrix.Username, config.Matrix.Password)
		if err != nil {
			log.Error().Err(err).Msg("Failed to auth")
		}
	} else {
		err = s.Login(config.Matrix.Username, config.Matrix.Password)
		if err != nil {
			log.Error().Err(err).Msg("Failed to auth")
		}
	}
	log.Info().Msg("Matrix access token updated")
}

func (s *service) StartListening(events chan *event.Event) (err error) {
	log.Info().Msg("Matrix event listener started")
	s.Syncer.OnEventType(event.StateMember, func(ctx context.Context, evt *event.Event) {
		eventTime := time.Unix(0, evt.Timestamp*int64(time.Millisecond))
		if eventTime.After(s.startTime) {
			log.Info().Msgf("Receive invite to room %v", evt.RoomID)
			if evt.GetStateKey() == s.Client.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
				_, err = s.Client.JoinRoomByID(s.Context, evt.RoomID)
				log.Debug().Msgf("Joined to room %v", evt.RoomID)
				if err != nil {
					log.Error().Err(err).Msg("Failed to join room")
				}
			}
		}
	})

	s.Syncer.OnEventType(event.EventMessage, func(ctx context.Context, evt *event.Event) {
		if evt.Sender == s.Client.UserID {
			return // Ignore messages sent by the bot
		}
		eventTime := time.Unix(0, evt.Timestamp*int64(time.Millisecond))
		if eventTime.After(s.startTime) {
			events <- evt
		}
	})

	go func() {
		for {
			log.Info().Msg("Matrix events sync started")
			err = s.Client.Sync()
			if err != nil {
				log.Error().Err(err).Msg("Sync failed")
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return
}

func isReply(evt *event.Event) bool {
	return evt.Content.AsMessage().RelatesTo != nil && evt.Content.AsMessage().RelatesTo.InReplyTo != nil
}

func (s *service) IsPrivateRoom(roomID id.RoomID) (bool, error) {
	members, err := s.Client.JoinedMembers(s.Context, roomID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get joined members")
		return false, err
	}

	return len(members.Joined) == 2, nil
}

func (s *service) GetRepliedEvent(evt *event.Event) (*event.Event, error) {

	if !isReply(evt) {
		return nil, errors.New("Not a reply")
	}

	return s.Client.GetEvent(s.Context, evt.RoomID, evt.Content.AsMessage().RelatesTo.InReplyTo.EventID)
}

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

func (s *service) SendMessage(roomID id.RoomID, text string) (*mautrix.RespSendEvent, error) {

	content := event.MessageEventContent{
		MsgType:       event.MsgText,
		Body:          text,
		Format:        event.FormatHTML,
		FormattedBody: fmt.Sprintf("<p>%s</p>", text),
	}
	return s.Client.SendMessageEvent(s.Context, roomID, event.EventMessage, content)
}

func (s *service) UploadAudio(base64Audio string) (id.ContentURIString, error) {
	audioBytes, err := base64.StdEncoding.DecodeString(base64Audio)
	if err != nil {
		log.Error().Err(err).Msg("Base64 decoding error")
	}

	uploadRequest := mautrix.ReqUploadMedia{
		ContentBytes:  audioBytes,
		ContentLength: int64(len(audioBytes)),
		ContentType:   "audio/webm",
		FileName:      uuid.New().String() + ".webm",
	}

	response, err := s.Client.UploadMedia(s.Context, uploadRequest)
	if err != nil {
		log.Error().Err(err).Msg("Uploading error")
	}
	log.Info().Msgf("Media upload response: %v", response.ContentURI)
	mxcURI := "mxc://" + response.ContentURI.Homeserver + "/" + response.ContentURI.FileID
	return id.ContentURIString(mxcURI), nil
}

func (s *service) SendMessageWithMedia(roomID, mxcURI, message string) error {

	jsonValue, _ := json.Marshal(map[string]interface{}{
		"msgtype":        "m.text",
		"body":           message,
		"format":         "org.matrix.custom.html",
		"formatted_body": fmt.Sprintf("%s\n<img src=\"%s\" alt=\"alt text\"/>", message, mxcURI),
	})

	request, err := http.NewRequest("POST", config.Matrix.HomeserverURL+"/_matrix/client/r0/rooms/"+roomID+"/send/m.room.message?access_token="+s.Client.AccessToken, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
