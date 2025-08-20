package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix"
	"github.com/tensved/bobrix/mxbot"
	"github.com/tensved/snet-matrix-framework/internal/app"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/snet"
)

func main() {
	log.Info().Msg("starting SNET Matrix Framework")

	a := app.New()
	log.Info().Msg("application initialized")

	// Authentication loop with graceful context handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().Msg("starting matrix authentication")
	a.MatrixClient.Auth()
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Debug().Msg("refreshing matrix authentication")
			a.MatrixClient.Auth()
		}
	}()

	log.Info().Msg("starting initial sync")
	a.Syncer.SyncOnce()
	log.Info().Msg("initial sync completed")

	engine := bobrix.NewEngine()
	log.Info().Msg("bobrix engine created")

	botCredentials := &mxbot.BotCredentials{
		Username:      config.Matrix.Username,
		Password:      config.Matrix.Password,
		HomeServerURL: config.Matrix.HomeserverURL,
		PickleKey:     []byte(config.Matrix.PickleKey),
	}
	log.Info().Str("username", config.Matrix.Username).Str("homeserver", config.Matrix.HomeserverURL).Msg("bot credentials prepared")

	snetBot, err := snet.NewSNETBot(botCredentials, a.MatrixClient, a.Ethereum, a.DB, a.GRPCManager, a.Syncer.FileDescriptors)
	if err != nil {
		log.Error().Err(err).Msg("failed to create snet bot")
		panic(err)
	}
	log.Info().Msg("SNET bot created successfully")

	engine.ConnectBot(snetBot)
	log.Info().Msg("bot connected to engine")

	go a.Syncer.Start(ctx)
	log.Info().Msg("syncer started")

	// Start the Fiber server in a goroutine with error handling
	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Msg("starting fiber server")
		a.Fiber.RegisterFiberRoutes()
		if err = a.Fiber.App.Listen("0.0.0.0:" + config.App.Port); err != nil {
			log.Error().Err(err).Msg("fiber server failed to start")
			serverErrors <- err
		}
	}()

	log.Info().Msg("starting bot engine")
	if err = engine.Run(ctx); err != nil {
		log.Error().Err(err).Msg("failed to run bot engine")
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Select loop to handle shutdown or server error
	select {
	case <-quit:
		log.Info().Msg("shutting down...")

	case err = <-serverErrors:
		log.Error().Err(err).Msg("server error, shutting down")
	}

	log.Info().Msg("stopping syncer")
	a.Syncer.Stop()

	cancel()

	log.Info().Msg("stopping bot engine")
	if err = engine.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("failed to stop bot engine")
	}
	log.Info().Msg("shutting down fiber server")
	if err = a.Fiber.App.Shutdown(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown fiber server")
	}

	log.Info().Msg("service shutdown complete")
}
