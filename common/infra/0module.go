package infra

import (
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-support/internal/infra")

func Module() fx.Option {
	return fx.Module(
		"internal.infra",
		fx.Provide(NewCacheDriver),
		fx.Provide(fx.Annotate(NewDCron, fx.ParamTags(``, `name:"serviceName"`))),
		fx.Provide(NewEtcdClient),
		fx.Provide(NewPrometheusRegistry),
		fx.Invoke(InvokePrometheusGoCollector),
		fx.Provide(NewRedisClient),
		fx.Decorate(InjectRedisObservability),
		fx.Provide(NewRedisRateLimiter),
	)
}
