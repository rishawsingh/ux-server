package middlewares

import (
	"github.com/remotestate/golang/internal"
	cors "github.com/rs/cors/wrapper/gin"
)

type CorsMiddleware struct {
	handler *internal.RequestHandler
	logger  *internal.Logger
	env     *internal.EnvConfig
}

func NewCorsMiddleware(handler *internal.RequestHandler, logger *internal.Logger, env *internal.EnvConfig) *CorsMiddleware {
	return &CorsMiddleware{
		handler: handler,
		logger:  logger,
		env:     env,
	}
}

func (m *CorsMiddleware) Setup() {
	m.logger.Info("Setting up cors middleware")
	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
	}))
}
