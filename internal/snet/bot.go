package snet

import (
	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix"
	"github.com/tensved/bobrix/contracts"
	"github.com/tensved/bobrix/mxbot"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/internal/matrix"
	"github.com/tensved/snet-matrix-framework/internal/syncer"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewSNETBot(credentials *mxbot.BotCredentials, matrix matrix.Service, eth blockchain.Ethereum, database db.Service, grpc *grpcmanager.GRPCClientManager, fileDescriptors map[string][]protoreflect.FileDescriptor) (*bobrix.Bobrix, error) {
	logger := log.With().
		Str("bot_name", "snet").
		Str("username", credentials.Username).
		Logger()

	bot, err := mxbot.NewDefaultBot("snet", credentials)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create bot")
		return nil, err
	}

	logger.Debug().Msg("bot created successfully")

	bot.AddCommand(mxbot.NewCommand(
		"info",
		func(c mxbot.CommandCtx) error {
			logger.Debug().
				Str("room", string(c.Event().RoomID)).
				Str("sender", c.Event().Sender.String()).
				Msg("info command received")

			info := syncer.GetSnetServicesInfo(fileDescriptors)
			logger.Debug().
				Str("info", info).
				Msg("snet services info generated")

			err := c.TextAnswer(info)
			if err != nil {
				logger.Error().Err(err).Msg("failed to send info answer")
			}
			return err
		}, mxbot.CommandConfig{
			Prefix: "!",
			Description: map[string]string{
				"en": "Info about snet services",
				"ru": "Информация о сервисах SNET",
			},
		}),
	)

	// Verbose event logging
	bot.AddEventHandler(
		mxbot.NewLoggerHandler("snet"),
	)

	// Add custom invite logging handler
	bot.AddEventHandler(
		mxbot.NewStateMemberHandler(func(ctx mxbot.Ctx) error {
			evt := ctx.Event()
			stateKeyStr := "nil"
			if evt.StateKey != nil {
				stateKeyStr = *evt.StateKey
			}

			logger.Debug().
				Str("room", string(evt.RoomID)).
				Str("sender", evt.Sender.String()).
				Str("state_key", stateKeyStr).
				Str("bot_full_name", bot.FullName()).
				Str("membership", string(evt.Content.AsMember().Membership)).
				Msg("membership event received")

			// Check if this is an invite for our bot
			if evt.Content.AsMember().Membership == "invite" && evt.StateKey != nil && *evt.StateKey == bot.FullName() {
				logger.Info().
					Str("room", string(evt.RoomID)).
					Str("inviter", evt.Sender.String()).
					Msg("bot invite detected")
			}
			return nil
		}, mxbot.FilterMembershipInvite()),
	)

	// Auto-join on invites with additional logs
	logger.Debug().Msg("registering auto-join room handler")
	bot.AddEventHandler(
		mxbot.AutoJoinRoomHandler(bot, mxbot.JoinRoomParams{
			PreJoinHook: func(ctx mxbot.Ctx) error {
				logger.Info().
					Str("room", string(ctx.Event().RoomID)).
					Str("inviter", ctx.Event().Sender.String()).
					Msg("joining room after invite")
				return nil
			},
			AfterJoinHook: func(ctx mxbot.Ctx) error {
				logger.Info().
					Str("room", string(ctx.Event().RoomID)).
					Msg("successfully joined room")
				return nil
			},
		}),
	)

	callStates := make(map[string]*CallState)

	bobr := bobrix.NewBobrix(bot)
	bobr.SetContractParser(Parser(matrix, eth, database, callStates, bobr, grpc))

	// Connect services to the bot from file descriptors.
	if len(fileDescriptors) > 0 && bobr != nil && database != nil && grpc != nil {
		logger.Info().
			Int("services_count", len(fileDescriptors)).
			Msg("connecting services to bot")
		createServices(bobr, fileDescriptors, eth, database, grpc)
	} else {
		logger.Warn().Msg("no services to connect or missing dependencies")
	}

	logger.Info().Msg("SNET bot initialization completed")
	return bobr, nil
}

func createServices(bobr *bobrix.Bobrix, fileDescriptors map[string][]protoreflect.FileDescriptor, eth blockchain.Ethereum, database db.Service, grpc *grpcmanager.GRPCClientManager) {
	logger := log.With().
		Int("total_services", len(fileDescriptors)).
		Logger()

	for snetIDOfService, descriptors := range fileDescriptors {
		logger.Debug().
			Str("snet_id", snetIDOfService).
			Int("descriptors_count", len(descriptors)).
			Msg("processing service descriptors")

		for _, descriptor := range descriptors {
			if descriptor == nil {
				logger.Warn().Msg("skipping nil descriptor")
				continue
			}

			logger.Debug().
				Str("descriptor_name", string(descriptor.FullName())).
				Msg("processing descriptor")

			services := descriptor.Services()
			if services == nil {
				logger.Warn().
					Str("descriptor", string(descriptor.FullName())).
					Msg("no services found in descriptor")
				continue
			}

			for i := range services.Len() {
				serviceDescriptor := services.Get(i)
				if serviceDescriptor == nil {
					logger.Warn().
						Int("index", i).
						Msg("skipping nil service descriptor")
					continue
				}

				serviceName := serviceDescriptor.Name()
				logger.Info().
					Str("snet_id", snetIDOfService).
					Str("service", string(serviceName)).
					Str("descriptor", string(descriptor.FullName())).
					Msg("connecting service")

				bobr.ConnectService(NewService(serviceDescriptor, string(descriptor.FullName()), snetIDOfService, string(serviceName), eth, database, grpc), func(ctx mxbot.Ctx, r *contracts.MethodResponse, _ any) {
					if r == nil {
						logger.Error().Msg("service returned nil response")
						_ = ctx.TextAnswer("Unexpected error")
						return
					}
					if r.Err != nil {
						logger.Error().
							Err(r.Err).
							Int("error_code", r.ErrCode).
							Msg("service handler error")
						_ = ctx.ErrorAnswer(r.Err.Error(), r.ErrCode)
						return
					}

					answer, ok := r.GetString("answer")
					if !ok || answer == "" {
						answer = "Unexpected error"
					}

					if err := ctx.TextAnswer(answer); err != nil {
						logger.Error().
							Err(err).
							Str("answer", answer).
							Msg("failed to send text answer")
					}
				})
			}
		}
	}

	logger.Info().Msg("service creation completed")
}
