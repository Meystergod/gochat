package main

import (
	"context"

	"github.com/Meystergod/gochat/internal/app"
	"github.com/Meystergod/gochat/internal/config"
	"github.com/Meystergod/gochat/pkg/logging"
)

func main() {
	logger := logging.NewDefaultLogger()

	cfg := config.GetConfig()
	if cfg == nil {
		logger.Fatal().Msg("reading config error")
	}

	loggerDeps := &logging.LoggerDeps{
		LogLevel: cfg.Log.LogLevel,
		LogFile:  cfg.Log.LogFile,
		LogSize:  cfg.Log.LogSize,
		LogAge:   cfg.Log.LogAge,
	}
	logger, err := logging.NewLogger(loggerDeps)
	if err != nil {
		logger.Fatal().Msgf("creating logger error: %s", err)
	}

	logger.Debug().Msg("initialized logger")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logger.WithContext(ctx)

	application, err := app.NewApplication(ctx, cfg)
	if err != nil {
		logger.Fatal().Msgf("creating application error: %s", err)
	}

	logger.Debug().Msg("created new application")

	if err := application.Run(ctx); err != nil {
		logger.Fatal().Msgf("running application error: %s", err)
	}

	logger.Debug().Msg("service done")
}
