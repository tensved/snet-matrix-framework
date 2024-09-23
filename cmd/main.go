package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix"
	"github.com/tensved/bobrix/mxbot"
	"github.com/tensved/snet-matrix-framework/internal/app"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/snet"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	a := app.New()

	// Authentication loop with graceful context handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.MatrixClient.Auth()
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			a.MatrixClient.Auth()
		}
	}()

	a.Syncer.SyncOnce()

	engine := bobrix.NewEngine()

	botCredentials := &mxbot.BotCredentials{
		Username:      config.Matrix.Username,
		Password:      config.Matrix.Password,
		HomeServerURL: config.Matrix.HomeserverURL,
	}

	snetBot, err := snet.NewSNETBot(botCredentials, a.MatrixClient, a.Ethereum, a.DB, a.GRPCManager, a.Syncer.FileDescriptors)
	if err != nil {
		log.Error().Err(err).Msg("failed to create snet bot")
	}

	engine.ConnectBot(snetBot)

	go a.Syncer.Start(ctx)

	// Start the Fiber server in a goroutine with error handling
	serverErrors := make(chan error, 1)
	go func() {
		a.Fiber.RegisterFiberRoutes()
		if err = a.Fiber.App.Listen(":" + config.App.Port); err != nil {
			serverErrors <- err
		}
	}()

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

	// Stop the Syncer
	a.Syncer.Stop()

	// Initiate shutdown
	cancel()

	// Stop the bot engine
	if err = engine.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("failed to stop bot engine")
	}
	// Stop the fiber server
	if err = a.Fiber.App.Shutdown(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown fiber server")
	}

	log.Info().Msg("service shutdown complete")
}
