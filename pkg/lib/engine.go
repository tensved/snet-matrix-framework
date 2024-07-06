package lib

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"log/slog"
	"matrix-ai-framework/internal/app"
	"matrix-ai-framework/internal/config"
	"time"
)

type IEngine interface {
	RegisterBot(bot MatrixBot)
	Run(ctx context.Context) error
}

type Engine struct {
	Bots   []MatrixBot
	Logger *slog.Logger
}

type EngineOpts func(*Engine)

func WithLogger(log *slog.Logger) EngineOpts {
	return func(e *Engine) {
		e.Logger = log
	}
}

func NewEngine(opts ...EngineOpts) *Engine {
	engine := &Engine{
		Logger: slog.Default(),
		Bots:   make([]MatrixBot, 0),
	}

	for _, opt := range opts {
		opt(engine)
	}
	return engine
}

func (e *Engine) RegisterBot(bot MatrixBot) {
	e.Bots = append(e.Bots, bot)
}

func (e *Engine) Run(ctx context.Context) error {

	e.Logger.Info("Starting engine...")

	for _, bot := range e.Bots {
		go func() {
			e.Logger.Info(fmt.Sprintf("Running bot: %s", ""))
			if err := bot.Run(); err != nil {
				e.Logger.Error(fmt.Sprintf("Failed to run bot: %s", err))
				return
			}
		}()
	}

	<-ctx.Done()

	return nil
}

func DefaultSNETEngine() *Engine {
	engine := NewEngine()

	a := app.New()

	go func() {
		a.Fiber.RegisterFiberRoutes()
		err := a.Fiber.Listen(":" + config.App.Port)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start fiber server")
		}
	}()

	a.MatrixClient.Auth()
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// Call the Auth function every time the ticker ticks.
				a.MatrixClient.Auth()
			}
		}
	}()

	a.Syncer.SyncOnce()
	go a.Syncer.Start()

	bot := NewSNETBot(a.MatrixClient, a.DB, a.Ethereum)

	// connect services to the bot from file descriptors
	if a.Syncer.FileDescriptors != nil {
		for snetIDOfService, descriptors := range a.Syncer.FileDescriptors {
			log.Debug().Msgf("service snet id: %s", snetIDOfService)
			if descriptors != nil {
				for _, descriptor := range descriptors {
					if descriptor != nil {
						log.Debug().Msgf("service descriptor name: %s", descriptor.FullName())
						services := descriptor.Services()
						if services != nil {
							for i := 0; i < services.Len(); i++ {
								if services.Get(i) != nil {
									serviceName := services.Get(i).Name()
									log.Debug().Msgf("service name: %s", services.Get(i).Name())
									serviceDescriptor := services.Get(i)
									botService := NewSNETService(serviceDescriptor, snetIDOfService, string(serviceName), a.Ethereum, a.DB, a.GRPCManager)
									bot.ConnectService(botService, BotServiceOpts{})
								}
							}
						}
					}
				}
			}
		}
	}

	engine.RegisterBot(bot)

	return engine
}
