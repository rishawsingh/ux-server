package utils

import "go.uber.org/fx"

// Module exports utils dependencies
var Module = fx.Options(
	fx.Provide(NewInMemScore),
)
