package internal

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/remotestate/golang/constants"

	"github.com/gin-gonic/gin"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(env *EnvConfig, logger *Logger, tracer *Tracer) *RequestHandler {
	if env.Environment == Production {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()
	engine.ForwardedByClientIP = true
	engine.Use(gin.Recovery())
	engine.Use(tracer.GinMiddleware())
	engine.Use(func(c *gin.Context) {
		traceID, spanID, found := tracer.GetTraceAndSpanID(c)
		if !found {
			traceID, spanID = uuid.NewString(), uuid.NewString()
		}
		// check if traceID and spanID is coming in request to service
		if c.GetHeader(constants.TraceID.String()) != "" && c.GetHeader(constants.SpanID.String()) != "" {
			traceID = c.GetHeader(constants.TraceID.String())
			spanID = c.GetHeader(constants.SpanID.String())
		}
		reqCtx := c.Request.Context()
		reqCtx = context.WithValue(reqCtx, constants.TraceID, traceID)
		reqCtx = context.WithValue(reqCtx, constants.SpanID, spanID)
		c.Request = c.Request.WithContext(reqCtx)
		c.Header("Content-Type", "application/json")
		c.Header(constants.TraceID.String(), traceID)
		c.Header(constants.SpanID.String(), spanID)
		c.Next()
	})
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Page not found",
		})
	})
	return &RequestHandler{Gin: engine}
}
