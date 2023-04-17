package services

import (
	"github.com/remotestate/golang/services/attachment"
	"github.com/remotestate/golang/services/group"
	"github.com/remotestate/golang/services/location"
	"github.com/remotestate/golang/services/product"
	"github.com/remotestate/golang/services/survey"
	"go.uber.org/fx"
)

// Module exports services present
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			product.NewProductService,
			fx.As(
				new(Product)),
		),
	),
	fx.Provide(
		fx.Annotate(
			survey.NewSurveyService,
			fx.As(
				new(Survey)),
		),
	),
	fx.Provide(
		fx.Annotate(
			group.NewGroupService,
			fx.As(
				new(Group)))),
	fx.Provide(
		fx.Annotate(
			attachment.NewAttachmentService,
			fx.As(
				new(Attachment)))),
	fx.Provide(
		fx.Annotate(
			location.NewLocationService,
			fx.As(
				new(Location)))),
)
