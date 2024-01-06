package services

import (
	"github.com/AH-dark/gravatar-with-qq-avatar/services/avatar"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(avatar.NewService),
	)
}
