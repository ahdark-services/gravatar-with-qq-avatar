package infra

import (
	"context"
	"fmt"

	"github.com/libi/dcron"
	"github.com/libi/dcron/driver"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-support/pkg/dcron_logger"
)

func NewDCron(ctx context.Context, serviceName string, vip *viper.Viper, etcdClient *clientv3.Client, lc fx.Lifecycle) *dcron.Dcron {
	ctx, span := tracer.Start(ctx, "infra.NewDCron")
	defer span.End()

	d := dcron.NewDcronWithOption(
		fmt.Sprintf("%s-%s", vip.GetString("name"), serviceName),
		driver.NewEtcdDriver(etcdClient),
		dcron.WithHashReplicas(10),
		dcron.WithLogger(dcron_logger.NewZapLogger(otelzap.L().Named("dcron"))),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx, span := tracer.Start(ctx, "infra.NewDCron.OnStart")
			defer span.End()

			d.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, span := tracer.Start(ctx, "infra.NewDCron.OnStop")
			defer span.End()

			d.Stop()

			return nil
		},
	})

	return d
}
