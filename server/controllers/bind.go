package controllers

import (
	"context"
	"github.com/AH-dark/gravatar-with-qq-avatar/server/middlewares"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/etag"
	"github.com/hertz-contrib/gzip"
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

	svr.Use(middlewares.RequestId())

	avatarRouter := svr.Group("/avatar")
	avatarRouter.Use(gzip.Gzip(gzip.BestCompression))
	avatarRouter.Use(etag.New())
	{
		avatarRouter.GET("", handlers.AvatarHandlers.GetAvatar)
		avatarRouter.GET("/:hash", handlers.AvatarHandlers.GetAvatar)
	}
}
