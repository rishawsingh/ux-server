package score

import (
	"github.com/remotestate/golang/api/controllers/score"
	"github.com/remotestate/golang/internal"
)

type Routes struct {
	logger     *internal.Logger
	handler    *internal.RequestHandler
	controller *score.Controller
}

func NewRoutes(logger *internal.Logger, handler *internal.RequestHandler, controller *score.Controller) *Routes {
	return &Routes{
		logger:     logger,
		handler:    handler,
		controller: controller,
	}
}

func (r *Routes) Setup() {
	r.logger.Info("setting up score routes")
	//api := r.handler.Gin.Group("/score/:surveyID")
	//api.GET("/calculate", r.controller.Calculate)
}
