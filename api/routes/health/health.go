package health

import (
	"github.com/remotestate/golang/api/controllers/health"
	"github.com/remotestate/golang/internal"
)

type Routes struct {
	logger     *internal.Logger
	handler    *internal.RequestHandler
	controller *health.Controller
}

func NewRoutes(logger *internal.Logger, handler *internal.RequestHandler, controller *health.Controller) *Routes {
	return &Routes{
		logger:     logger,
		handler:    handler,
		controller: controller,
	}
}

func (r *Routes) Setup() {
	r.logger.Info("setting up health routes")
	api := r.handler.Gin.Group("/health")
	api.GET("/", r.controller.Hello)
}
