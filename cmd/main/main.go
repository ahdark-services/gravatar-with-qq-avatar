package main

import (
	"context"
	"github.com/AH-dark/gravatar-with-qq-avatar/entry"
	"github.com/AH-dark/gravatar-with-qq-avatar/server"
	"go.uber.org/fx"
)

var (
	ctx = context.Background()
	app *fx.App
)

func init() {
	app = fx.New(
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fx.Supply(fx.Annotate("main", fx.ResultTags(`name:"serviceName"`))),
		entry.AppEntries(),
		server.Module(),
	)
}

func main() {
	app.Run()
}
