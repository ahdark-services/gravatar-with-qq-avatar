package database

import (
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-support/database/dal"
	"github.com/AH-dark/gravatar-with-qq-support/database/instances"
)

func Module() fx.Option {
	return fx.Module("database",
		fx.Provide(instances.NewClusterConfig),
		fx.Provide(instances.NewSession),

		fx.Provide(dal.NewMD5QQMapping),
	)
}
