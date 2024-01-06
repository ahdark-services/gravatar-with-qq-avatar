package server

import (
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-avatar/server/controllers"
	"github.com/AH-dark/gravatar-with-qq-avatar/server/controllers/avatar"
)

func Module() fx.Option {
	return fx.Module("server",
		fx.Provide(NewServer),
		fx.Invoke(RunServer),

		fx.Provide(avatar.NewHandlers),
		fx.Invoke(controllers.BindControllers),
	)
}
