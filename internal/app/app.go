package app

import (
	"context"
	"fmt"
	"github.com/Meystergod/gochat/internal/apperror"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/Meystergod/gochat/internal/config"
	"github.com/Meystergod/gochat/internal/controller"
	"github.com/Meystergod/gochat/internal/delivery/http/v1/httpecho"
	"github.com/Meystergod/gochat/internal/repository/repository_user/mongodb"
	"github.com/Meystergod/gochat/internal/usecase/usecase_user"
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

	a.httpServer.Server().HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appError *apperror.AppError
		if errors.As(err, &appError) {
			switch {
			case errors.Is(appError.Err, apperror.ErrorDecode),
				errors.Is(appError.Err, apperror.ErrorCreateOne),
				errors.Is(appError.Err, apperror.ErrorConvert),
				errors.Is(appError.Err, apperror.ErrorGetOne),
				errors.Is(appError.Err, apperror.ErrorGetAll),
				errors.Is(appError.Err, apperror.ErrorUpdateOne),
				errors.Is(appError.Err, apperror.ErrorDeleteOne),
				errors.Is(appError.Err, apperror.ErrorConvertModel):

				appError.ErrorMessage = appError.Error()
				if jsonError := c.JSON(http.StatusInternalServerError, &appError); jsonError != nil {
					logger.Error().Msgf("failed to create json response: %s", jsonError.Error())
				}
				return
			case errors.Is(appError.Err, apperror.ErrorValidatePayload):
				appError.ErrorMessage = appError.Error()
				if jsonError := c.JSON(http.StatusBadRequest, &appError); jsonError != nil {
					logger.Error().Msgf("failed to validate json response: %s", jsonError.Error())
				}
				return
			}
		}
		a.httpServer.Server().DefaultHTTPErrorHandler(err, c)
	}

	a.httpServer.Server().Validator = utils.NewValidator()

	userRepository := repository_user.NewUserRepository(a.db, utils.CollNameUser)
	userUsecase := usecase_user.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)

	httpecho.SetUserApiRoutes(a.httpServer.Server(), userController)
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

func (a *Application) migrate(_ context.Context) error {
	return nil
}
