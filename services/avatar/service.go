package avatar

import (
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
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
	fx.In
	MD5QQMappingRepo dal.MD5QQMappingRepo
}

func NewService(s service) Service {
	return &s
}

func (s *service) downloadAvatar(ctx context.Context, url string) (io.ReadCloser, *req.Response, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.downloadAvatar")
	defer span.End()

	resp, err := req.NewRequest().SetContext(ctx).Get(url)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("download avatar failed", zap.Error(err))
		return nil, resp, err
	}

	if resp.StatusCode != 200 || resp.IsErrorState() {
		otelzap.L().Ctx(ctx).Error("download avatar failed")
		return nil, resp, errors.New(http.StatusText(resp.StatusCode))
	}

	if resp.IsSuccessState() {
		return resp.Body, resp, nil
	}

	return nil, resp, fmt.Errorf("download avatar failed, status code: %d", resp.StatusCode)
}

func (s *service) GetAvatar(ctx context.Context, hash string, args GetAvatarArgs) ([]byte, time.Time, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.GetAvatar")
	defer span.End()

	var lastModified time.Time

	// get qq avatar
	img, err := s.getQQAvatar(ctx, hash, args)
	if err != nil {
		res, err := s.getGravatar(ctx, hash, args)
		if errors.Is(err, errors.New("404 Not Found")) {
			return nil, time.Time{}, nil
		} else if err != nil {
			otelzap.L().Ctx(ctx).Warn("get gravatar failed", zap.Error(err))
			return nil, time.Time{}, err
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
