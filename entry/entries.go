package entry

import (
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-support/common/config"
	"github.com/AH-dark/gravatar-with-qq-support/common/infra"
	"github.com/AH-dark/gravatar-with-qq-support/common/logging"
	"github.com/AH-dark/gravatar-with-qq-support/common/observability"
	"github.com/AH-dark/gravatar-with-qq-support/common/settings"
)

func AppEntries() fx.Option {
	return fx.Options(
		config.Module(),
		logging.Module(),
		fx.WithLogger(logging.FxLogger),
		observability.Module(),
		settings.Module(),
		infra.Module(),
	)
}
