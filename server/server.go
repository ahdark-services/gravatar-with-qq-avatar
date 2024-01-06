package server

import (
	"context"
	"fmt"
	"go.uber.org/fx"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	promclient "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"

	hertzloggerzap "github.com/AH-dark/gravatar-with-qq-avatar/pkg/hertzloggerzap"
	hertzprometheus "github.com/AH-dark/gravatar-with-qq-avatar/pkg/hertzprometheus"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/server")

func NewServer(ctx context.Context, vip *viper.Viper, promRegistry *promclient.Registry) *server.Hertz {
	ctx, span := tracer.Start(ctx, "server.NewServer")
	defer span.End()

	hlog.SetLogger(hertzloggerzap.NewLoggerWithZapLogger(otelzap.L().Named("hertz")))

	traceOption, cfg := hertztracing.NewServerTracer()
	svr := server.Default(
		traceOption,
		server.WithNetwork(vip.GetString("server.network")),
		server.WithHostPorts(fmt.Sprintf("%s:%d", vip.GetString("server.host"), vip.GetUint16("server.port"))),
		server.WithHandleMethodNotAllowed(true), server.WithTracer(hertzprometheus.NewServerTracer(
			"",
			"",
			hertzprometheus.WithRegistry(promRegistry),
			hertzprometheus.WithEnableGoCollector(false),
			hertzprometheus.WithDisableServer(true),
		)),
		server.WithTransport(netpoll.NewTransporter),
	)
	svr.Use(hertztracing.ServerMiddleware(cfg))

	return svr
}

func RunServer(ctx context.Context, svr *server.Hertz, lc fx.Lifecycle) error {
	ctx, span := tracer.Start(ctx, "server.RunServer")
	defer span.End()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go svr.Spin()

			return nil
		},
	})

	return nil
}
