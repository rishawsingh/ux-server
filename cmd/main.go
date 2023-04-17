package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/remotestate/golang/utils"

	"github.com/remotestate/golang/services"

	"github.com/remotestate/golang/api/controllers"
	"github.com/remotestate/golang/api/middlewares"
	"github.com/remotestate/golang/api/routes"
	"github.com/remotestate/golang/internal"
	"go.uber.org/fx"
)

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func startServer(
	middleware *middlewares.Middlewares,
	env *internal.EnvConfig,
	db *internal.Database,
	router *internal.RequestHandler,
	route *routes.Routes,
	logger *internal.Logger,
	//cron *crons.Cron,
	tracer *internal.Tracer,
	lifecycle fx.Lifecycle) {
	middleware.Setup()
	route.Setup()
	//cron.Setup()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", env.ServerPort),
		Handler:           router.Gin,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			migrationErr := db.MigrateUp("./migration")
			if migrationErr != nil {
				logger.With(migrationErr).Panic("failed to run database migration")
			} else {
				logger.Info("database migration done")
			}
			go func(logger *internal.Logger, srv *http.Server) {
				logger.Infof("Running server at port %s", env.ServerPort)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.With(err).Fatal("failed to start server")
				}
			}(logger, srv)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if dbCloseErr := db.Close(); dbCloseErr != nil {
				return dbCloseErr
			}
			if serverCloseErr := srv.Shutdown(ctx); serverCloseErr != nil {
				return serverCloseErr
			}
			if tracerCloseErr := tracer.Shutdown(ctx); tracerCloseErr != nil {
				return tracerCloseErr
			}
			logger.Info("Server exiting")
			return nil
		},
	})
}

func main() {
	var CommonModules = fx.Options(
		controllers.Module,
		middlewares.Module,
		routes.Module,
		internal.Module,
		services.Module,
		utils.Module,
		//crons.Module,
	)
	logger := internal.GetLogger()
	opts := fx.Options(fx.WithLogger(logger.GetFxLogger))
	app := fx.New(CommonModules, opts, fx.Invoke(startServer))
	if app.Err() != nil {
		logger.With(app.Err()).Panic("failed to start the app")
	}
	app.Run()
}
