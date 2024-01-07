package avatar

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type GetGravatarResult struct {
	Avatar       image.Image
	LastModified time.Time
}

func (s *service) getGravatar(ctx context.Context, hash string, args GetAvatarArgs) (GetGravatarResult, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.getGravatar")
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
		return GetGravatarResult{}, err
	}
	u.RawQuery = values.Encode()

	stream, resp, err := s.downloadAvatar(ctx, u.String())
	if err != nil {
		otelzap.L().Ctx(ctx).Error("download gravatar failed", zap.Error(err))
		return GetGravatarResult{}, err
	}

	img, err := jpeg.Decode(stream)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("failed to decode gravatar", zap.Error(err))
		return GetGravatarResult{}, err
	}

	if err := stream.Close(); err != nil {
		otelzap.L().Ctx(ctx).Warn("close stream failed", zap.Error(err))
	}

	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		otelzap.L().Ctx(ctx).Warn("parse last modified failed", zap.Error(err))
		lastModified = time.Time{}
	}

	return GetGravatarResult{
		Avatar:       img,
		LastModified: lastModified,
	}, nil
}
