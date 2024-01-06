package avatar

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"

	"github.com/AH-dark/gravatar-with-qq-avatar/services/avatar"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/internal/server/controllers/avatar")

type Handlers interface {
	GetAvatar(ctx context.Context, c *app.RequestContext)
}

type handlers struct {
	fx.In
	AvatarService avatar.Service
}

func NewHandlers(h handlers) Handlers {
	return &h
}
