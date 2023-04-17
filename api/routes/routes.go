package routes

import (
	"github.com/remotestate/golang/api/routes/group"
	"github.com/remotestate/golang/api/routes/health"
	"github.com/remotestate/golang/api/routes/score"
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(health.NewRoutes),
	fx.Provide(score.NewRoutes),
	fx.Provide(group.NewRoutes),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	healthRoutes *health.Routes,
	scoreRoutes *score.Routes,
	groupRoutes *score.Routes,
) *Routes {
	return &Routes{
		healthRoutes,
		scoreRoutes,
		groupRoutes,
	}
}

// Setup all the route
func (r *Routes) Setup() {
	for _, route := range *r {
		route.Setup()
	}
}
