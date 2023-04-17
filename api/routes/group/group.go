package group

import (
	"github.com/remotestate/golang/api/controllers/group"
	"github.com/remotestate/golang/internal"
)

type Routes struct {
	logger     *internal.Logger
	handler    *internal.RequestHandler
	controller *group.Controller
}

func NewRoutes(logger *internal.Logger, handler *internal.RequestHandler, controller *group.Controller) *Routes {
	return &Routes{
		logger:     logger,
		handler:    handler,
		controller: controller,
	}
}

func (r *Routes) Setup() {
	r.logger.Info("setting up group routes")
	api := r.handler.Gin.Group("/group")

	api.POST("/", r.controller.CreateGroup)
}
