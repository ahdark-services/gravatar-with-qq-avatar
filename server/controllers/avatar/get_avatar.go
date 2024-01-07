package avatar

import (
	"context"
	"net/http"
	"strings"
	"time"

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
	Force   string `query:"f"`
}

func (h *handlers) GetAvatar(ctx context.Context, c *app.RequestContext) {
	ctx, span := tracer.Start(ctx, "server.controllers.avatarData.GetAvatar")
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
		ForceDefault: req.Force == "y",
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

	avatarData, lastModified, err := h.AvatarService.GetAvatar(ctx, req.Hash, args)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("get avatar data failed", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if avatarData == nil {
		c.NotFound()
		return
	}

	c.Data(http.StatusOK, lo.If(args.EnableWebp, "image/webp").Else("image/png"), avatarData)
	if !lastModified.IsZero() {
		c.Header("Last-Modified", lastModified.UTC().Format(time.RFC1123))
	}
	c.Header("Cache-Control", "public, max-age=86400, immutable")
	c.Header("Expires", time.Now().Add(86400*time.Second).UTC().Format(time.RFC1123))
	c.Header("X-Content-Type-Options", "nosniff")
}
