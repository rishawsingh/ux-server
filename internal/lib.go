package internal

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewEnvConfig),
	fx.Provide(GetLogger),
	fx.Provide(NewTracer),
	fx.Provide(NewDatabase),
	fx.Provide(NewTransactor),
	fx.Provide(NewRequestHandler),
)
