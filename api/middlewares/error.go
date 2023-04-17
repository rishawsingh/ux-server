package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/remotestate/golang/constants"
	"github.com/remotestate/golang/internal"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ErrorMiddleware struct {
	handler *internal.RequestHandler
	logger  *internal.Logger
	env     *internal.EnvConfig
}

func NewErrorHandler(handler *internal.RequestHandler, logger *internal.Logger, env *internal.EnvConfig) *ErrorMiddleware {
	return &ErrorMiddleware{
		handler: handler,
		logger:  logger,
		env:     env,
	}
}

func (m *ErrorMiddleware) Setup() {
	m.logger.Info("Setting up error middleware")
	m.handler.Gin.Use(func(c *gin.Context) {
		c.Next()
		hasErr := len(c.Errors)
		if hasErr > 0 {
			span := trace.SpanFromContext(c.Request.Context())
			if span.IsRecording() {
				span.SetStatus(codes.Error, "failed to handle request")
			}
			traceID, _ := c.Request.Context().Value(constants.TraceID).(string)
			spanID, _ := c.Request.Context().Value(constants.SpanID).(string)
			for i := range c.Errors {
				if span.IsRecording() {
					span.RecordError(c.Errors[i])
				}
				m.logger.WithCtx(c.Request.Context()).
					With("path", c.FullPath()).
					With("response status", c.Writer.Status()).
					With("error", c.Errors[i].Error()).
					Error("handler returned an error")
			}
			if m.env.Environment == internal.Production && c.Writer.Status() > 499 {
				c.JSON(c.Writer.Status(), gin.H{
					"error":                    "internal server error occurred",
					constants.TraceID.String(): traceID,
					constants.SpanID.String():  spanID,
				})
			} else {
				c.JSON(c.Writer.Status(), c.Errors.JSON())
			}
		}
	})
}
