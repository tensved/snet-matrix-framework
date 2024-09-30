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
	bot, err := mxbot.NewDefaultBot("snet", credentials)
	if err != nil {
		return nil, err
	}

	bot.AddCommand(mxbot.NewCommand(
		"info",
		func(c mxbot.CommandCtx) error {
			return c.Answer(syncer.GetSnetServicesInfo(fileDescriptors))
		}, mxbot.CommandConfig{
			Prefix:      "!",
			Description: "Info about snet services",
		}),
	)

	bot.AddEventHandler(
		mxbot.AutoJoinRoomHandler(bot),
	)

	callStates := make(map[string]*CallState)

	bobr := bobrix.NewBobrix(bot)

	bobr.SetContractParser(Parser(matrix, eth, database, callStates, bobr, grpc))

	// Connect services to the bot from file descriptors.
	if len(fileDescriptors) > 0 && bobr != nil && database != nil && grpc != nil {
		createServices(bobr, fileDescriptors, eth, database, grpc)
	}

	return bobr, nil
}

func createServices(bobr *bobrix.Bobrix, fileDescriptors map[string][]protoreflect.FileDescriptor, eth blockchain.Ethereum, database db.Service, grpc *grpcmanager.GRPCClientManager) {
	for snetIDOfService, descriptors := range fileDescriptors {
		log.Debug().Msgf("service snet id: %s", snetIDOfService)
		for _, descriptor := range descriptors {
			if descriptor == nil {
				continue
			}
			log.Debug().Msgf("service descriptor name: %s", descriptor.FullName())
			services := descriptor.Services()
			if services == nil {
				continue
			}
			for i := range services.Len() {
				serviceDescriptor := services.Get(i)
				if serviceDescriptor == nil {
					continue
				}
				serviceName := serviceDescriptor.Name()
				bobr.ConnectService(NewService(serviceDescriptor, string(descriptor.FullName()), snetIDOfService, string(serviceName), eth, database, grpc), func(ctx mxbot.Ctx, r *contracts.MethodResponse) {
					answer, ok := r.Data["answer"].(string)
					if !ok {
						answer = "Unexpected error"
					}

					if err := ctx.Answer(answer); err != nil {
						log.Error().Err(err).Msg("failed to answer")
					}
				})
			}
		}
	}
}
