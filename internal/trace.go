package internal

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	logger   *Logger
	enabled  bool
	provider *sdktrace.TracerProvider
	env      *EnvConfig
}

func NewTracer(logger *Logger, env *EnvConfig) (*Tracer, error) {
	tracer := &Tracer{
		logger: logger,
		env:    env,
	}
	if !env.Trace {
		logger.Info("skipping tracer, not enabled")
		return tracer, nil
	}
	tp, jaegerErr := SignozProvider(env, logger)
	if jaegerErr != nil {
		return nil, jaegerErr
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer.provider = tp
	tracer.enabled = true
	return tracer, nil
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.enabled {
		return t.provider.Shutdown(ctx)
	}
	return nil
}

func (t *Tracer) WrapDB(driverName, dbName, connString string) (*sql.DB, error) {
	return otelsql.Open(
		driverName,
		connString,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(dbName),
		otelsql.WithTracerProvider(t.provider))
}

func (t *Tracer) GinMiddleware() gin.HandlerFunc {
	if t.enabled {
		return otelgin.Middleware(t.env.ServiceName, otelgin.WithTracerProvider(t.provider))
	}
	return otelgin.Middleware(t.env.ServiceName)
}

func (t *Tracer) GetTraceAndSpanID(c *gin.Context) (traceID, spanID string, found bool) {
	if oteltrace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
		traceID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
		spanID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()
		return traceID, spanID, true
	}
	return "", "", false
}

func (t *Tracer) Span(ctx context.Context) (context.Context, oteltrace.Span) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	callerName := lo.Ternary[string](ok && details != nil, details.Name(), "unknown-caller")
	ctx, span := otel.GetTracerProvider().Tracer(t.env.ServiceName).Start(ctx, callerName)
	return ctx, span
}

func JaegerProvider(env *EnvConfig, logger *Logger) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(env.TracerEndpoint), jaeger.WithAgentPort(env.TracerPort)))
	if err != nil {
		logger.With(err).Error("failed to initialise jaeger tracer")
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(env.ServiceName),
			semconv.DeploymentEnvironmentKey.String(string(env.Environment)),
			semconv.ServiceVersionKey.String(env.BuildNumber()),
		)),
	)
	return provider, nil
}

func SignozProvider(env *EnvConfig, logger *Logger) (*sdktrace.TracerProvider, error) {
	exp, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", env.TracerEndpoint, env.TracerPort)),
		),
	)
	if err != nil {
		logger.With(err).Error("failed to initialise signoz tracer")
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(env.ServiceName),
			semconv.DeploymentEnvironmentKey.String(string(env.Environment)),
			semconv.ServiceVersionKey.String(env.BuildNumber()),
		)),
	)
	return provider, nil
}
