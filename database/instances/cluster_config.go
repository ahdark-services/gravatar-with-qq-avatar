package instances

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/database/instances")

func NewClusterConfig(ctx context.Context, vip *viper.Viper) *gocql.ClusterConfig {
	ctx, span := tracer.Start(ctx, "database.instances.NewClusterConfig")
	defer span.End()

	cluster := gocql.NewCluster(
		vip.GetStringSlice("cassandra.hosts")...,
	)

	cluster.Keyspace = vip.GetString("cassandra.keyspace")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Compressor = &gocql.SnappyCompressor{}
	cluster.Logger = &StdLogger{logger: otelzap.L().Named("cassandra")}

	if vip.GetString("cassandra.username") != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: vip.GetString("cassandra.username"),
			Password: vip.GetString("cassandra.password"),
		}
	}

	return cluster
}
