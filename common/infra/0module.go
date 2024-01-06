package infra

import (
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/internal/infra")

func Module() fx.Option {
	return fx.Module(
		"internal.infra",
		fx.Provide(NewCacheDriver),
		fx.Provide(NewPrometheusRegistry),
		fx.Invoke(InvokePrometheusGoCollector),
		fx.Provide(NewRedisClient),
		fx.Decorate(InjectRedisObservability),
		fx.Provide(NewRedisRateLimiter),
	)
}
