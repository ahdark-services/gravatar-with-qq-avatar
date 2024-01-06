package controllers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-avatar/server/controllers/avatar"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/internal/server/controllers")

type HandlerGroup struct {
	fx.In
	AvatarHandlers avatar.Handlers
}

func BindControllers(ctx context.Context, svr *server.Hertz, handlers HandlerGroup) {
	ctx, span := tracer.Start(ctx, "server.controllers.BindControllers")
	defer span.End()

	svr.GET("/avatar", handlers.AvatarHandlers.GetAvatar)
	svr.GET("/avatar/:hash", handlers.AvatarHandlers.GetAvatar)
}
