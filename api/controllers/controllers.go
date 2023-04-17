package controllers

import (
	"github.com/remotestate/golang/api/controllers/group"
	"github.com/remotestate/golang/api/controllers/health"
	"github.com/remotestate/golang/api/controllers/score"
	"go.uber.org/fx"
)

// Module exported for initializing controllers
var Module = fx.Options(
	fx.Provide(health.NewController),
	fx.Provide(score.NewController),
	fx.Provide(group.NewController),
)
