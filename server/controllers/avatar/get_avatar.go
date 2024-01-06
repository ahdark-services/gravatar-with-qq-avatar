package avatar

import (
	"context"
	"strings"

	"github.com/AH-dark/bytestring"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/samber/lo"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/AH-dark/gravatar-with-qq-avatar/services/avatar"
)

type GetAvatarRequest struct {
	Hash string `path:"hash"`

	Size    int64  `query:"s"`
	Default string `query:"d"`
	Rating  string `query:"r"`
	Force   bool   `query:"f"`
}

func (h *handlers) GetAvatar(ctx context.Context, c *app.RequestContext) {
	ctx, span := tracer.Start(ctx, "server.controllers.avatar.GetAvatar")
	defer span.End()

	var req GetAvatarRequest
	if err := c.Bind(&req); err != nil {
		otelzap.L().Ctx(ctx).Error("bind request failed", zap.Error(err))
		c.AbortWithStatus(400)
		return
	}

	args := avatar.GetAvatarArgs{
		Size:         req.Size,
		Default:      req.Default,
		ForceDefault: req.Force,
		Rating:       req.Rating,
		EnableWebp: lo.SomeBy(strings.Split(bytestring.BytesToString(c.GetHeader("Accept")), ","), func(item string) bool {
			return strings.HasPrefix(item, "image/webp")
		}),
	}

	if args.Size <= 0 {
		args.Size = 80
	}

	if args.EnableWebp {
		c.Header("Content-Type", "image/webp")
	} else {
		c.Header("Content-Type", "image/png")
	}

	qqAvatar, err := h.AvatarService.GetQQAvatar(ctx, req.Hash, args)
	if err != nil {
		otelzap.L().Ctx(ctx).Warn("get qq avatar failed", zap.Error(err))
	} else {
		c.Status(200)
		c.SetBodyStream(qqAvatar, -1)
		return
	}

	gravatar, err := h.AvatarService.GetGravatar(ctx, req.Hash, args)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("get gravatar failed", zap.Error(err))
		c.AbortWithStatus(500)
		return
	}
	defer gravatar.Close()

	c.Status(200)
	c.SetBodyStream(gravatar, -1)
}
