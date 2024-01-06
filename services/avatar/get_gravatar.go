package avatar

import (
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"io"
	"net/url"
	"strconv"
)

func (s *service) GetGravatar(ctx context.Context, hash string, args GetAvatarArgs) (io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.GetGravatar")
	defer span.End()

	values := url.Values{}
	values.Set("d", args.Default)
	values.Set("s", strconv.FormatInt(args.Size, 10))
	values.Set("r", args.Rating)

	// force default
	if args.ForceDefault {
		hash = ""
	}

	// get gravatar
	u, err := url.Parse(fmt.Sprintf("https://www.gravatar.com/avatar/%s", hash))
	if err != nil {
		otelzap.L().Ctx(ctx).Error("parse gravatar url failed", zap.Error(err))
		return nil, err
	}
	u.RawQuery = values.Encode()

	return s.downloadAvatar(ctx, u.String())
}
