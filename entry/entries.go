package entry

import (
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-avatar/common/config"
	"github.com/AH-dark/gravatar-with-qq-avatar/common/infra"
	"github.com/AH-dark/gravatar-with-qq-avatar/common/logging"
	"github.com/AH-dark/gravatar-with-qq-avatar/common/observability"
	"github.com/AH-dark/gravatar-with-qq-avatar/database"
	"github.com/AH-dark/gravatar-with-qq-avatar/services"
)

func AppEntries() fx.Option {
	return fx.Options(
		config.Module(),
		logging.Module(),
		fx.WithLogger(logging.FxLogger),
		observability.Module(),
		infra.Module(),
		database.Module(),
		services.Module(),
	)
}
