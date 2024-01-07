package avatar

import (
	"context"
	"github.com/samber/lo"
	"image/png"
	"time"

	"github.com/cloudwego/hertz/pkg/common/bytebufferpool"
	"github.com/imroc/req/v3"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AH-dark/gravatar-with-qq-avatar/database/dal"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/services/avatar")

type GetAvatarArgs struct {
	Size         int64
	Default      string
	ForceDefault bool
	Rating       string

	EnableWebp bool
}

type Service interface {
	GetAvatar(ctx context.Context, hash string, args GetAvatarArgs) ([]byte, time.Time, error)
}

type service struct {
	fx.In            `ignore-unexported:"true"`
	MD5QQMappingRepo dal.MD5QQMappingRepo

	qqAvatarClient *req.Client
	gravatarClient *req.Client
}

func NewService(s service) Service {
	s.qqAvatarClient = initQQClient()
	s.gravatarClient = initGravatarClient()

	return &s
}

func (s *service) GetAvatar(ctx context.Context, hash string, args GetAvatarArgs) ([]byte, time.Time, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.GetAvatar")
	defer span.End()

	var lastModified time.Time

	// get qq avatar
	img, err := s.getQQAvatar(ctx, hash, args)
	if err != nil {
		res, err := s.getGravatar(ctx, hash, args)
		if err != nil {
			otelzap.L().Ctx(ctx).Warn("get gravatar failed", zap.Error(err))
			return nil, time.Time{}, err
		}

		if lo.IsEmpty(res) {
			return nil, time.Time{}, nil
		}

		img = res.Avatar
		lastModified = res.LastModified
	}

	// encode image
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	if args.EnableWebp {
		if err := webp.Encode(b, img, nil); err != nil {
			otelzap.L().Ctx(ctx).Error("encode webp failed", zap.Error(err))
			return nil, time.Time{}, err
		}
	} else {
		if err := png.Encode(b, img); err != nil {
			otelzap.L().Ctx(ctx).Error("encode png failed", zap.Error(err))
			return nil, time.Time{}, err
		}
	}

	return b.Bytes(), lastModified, nil
}
