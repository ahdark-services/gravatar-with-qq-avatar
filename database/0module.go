package database

import (
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-avatar/database/dal"
	"github.com/AH-dark/gravatar-with-qq-avatar/database/instances"
)

func Module() fx.Option {
	return fx.Module("database",
		fx.Provide(instances.NewClusterConfig),
		fx.Provide(instances.NewSession),

		fx.Provide(dal.NewMD5QQMapping),
	)
}
