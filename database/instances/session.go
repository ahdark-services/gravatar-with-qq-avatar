package instances

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewSession(ctx context.Context, cluster *gocql.ClusterConfig, lc fx.Lifecycle) (*gocqlx.Session, error) {
	ctx, span := tracer.Start(ctx, "database.instances.NewSession")
	defer span.End()

	sess, err := otelgocql.NewSessionWithTracing(ctx, cluster)
	if err != nil {
		otelzap.L().Ctx(ctx).Panic("create session failed", zap.Error(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			otelzap.L().Ctx(ctx).Info("closing cassandra connection")
			sess.Close()
			return nil
		},
	})

	cqlxSession := gocqlx.NewSession(sess)
	return &cqlxSession, nil
}
