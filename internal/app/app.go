package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Meystergod/gochat/internal/config"
	"github.com/Meystergod/gochat/internal/controller"
	"github.com/Meystergod/gochat/internal/delivery/http/v1/httpecho"
	"github.com/Meystergod/gochat/internal/repository/mongorepo"
	"github.com/Meystergod/gochat/internal/usecase"
	"github.com/Meystergod/gochat/internal/utils"
	"github.com/Meystergod/gochat/pkg/client"
	"github.com/Meystergod/gochat/pkg/httpserver"
	"github.com/Meystergod/gochat/pkg/ossignal"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	cfg        *config.Config
	httpServer *httpserver.Server
	db         *mongo.Database
}

func NewApplication(ctx context.Context, cfg *config.Config) (*Application, error) {
	dbConfig := client.NewMongoConfig(
		cfg.Database.Auth,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := client.NewMongoDatabase(ctx, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "connecting database")
	}

	return &Application{
		cfg:        cfg,
		httpServer: nil,
		db:         db,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	runner, ctx := errgroup.WithContext(ctx)

	runner.Go(func() error {
		if err := a.startHTTP(ctx); err != nil {
			return errors.Wrap(err, "listening and starting http api")
		}

		return nil
	})

	runner.Go(func() error {
		if err := ossignal.DefaultSignalWaiter(ctx); err != nil {
			return errors.Wrap(err, "waiting os signal")
		}

		return nil
	})

	runner.Go(func() error {
		<-ctx.Done()

		ctxSignal, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		addr := fmt.Sprintf("%s:%s", a.cfg.HTTPServer.Host, a.cfg.HTTPServer.Port)
		logger.Info().Str("addr", addr).Msg("shutdown http server")

		if err := a.shutdownHTTP(ctxSignal); err != nil {
			logger.Error().Err(err).Msg("shutdown http server")
		}

		return nil
	})

	if err := runner.Wait(); err != nil {
		switch {
		case ossignal.IsExitSignal(err):
			logger.Info().Msg("exited by exit signal")
		default:
			return errors.Wrap(err, "exiting with error")
		}
	}

	return nil
}

func (a *Application) startHTTP(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	httpServerDeps := &httpserver.ServerDeps{
		Host: a.cfg.HTTPServer.Host,
		Port: a.cfg.HTTPServer.Port,
	}

	a.httpServer = httpserver.NewDefaultServer(httpServerDeps)
	logger.Debug().Msg("created new http server")

	userRepository := mongorepo.NewUserRepository(a.db, utils.CollNameUser)
	userController := controller.NewUserController(userRepository)
	userUsecase := usecase.NewUserUsecase(userController)

	httpecho.SetUserApiRoutes(a.httpServer.Server(), userUsecase)
	logger.Debug().Msg("set api routes for user")

	addr := fmt.Sprintf("%s:%s", a.cfg.HTTPServer.Host, a.cfg.HTTPServer.Port)
	logger.Info().Str("addr", addr).Msg("listen and serve http api")

	err := a.httpServer.Start()
	if err != nil {
		return errors.Wrap(err, "starting http server error")
	}

	return nil
}

func (a *Application) shutdownHTTP(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
