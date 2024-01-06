package avatar

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"

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
	GetGravatar(ctx context.Context, hash string, args GetAvatarArgs) (io.ReadCloser, error)
	GetQQAvatar(ctx context.Context, hash string, args GetAvatarArgs) (io.ReadCloser, error)
}

type service struct {
	fx.In
	MD5QQMappingRepo dal.MD5QQMappingRepo
}

func NewService(s service) Service {
	return &s
}

func (s *service) downloadAvatar(ctx context.Context, url string) (io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.downloadAvatar")
	defer span.End()

	resp, err := req.NewRequest().SetContext(ctx).Get(url)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("download avatar failed", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != 200 || resp.IsErrorState() {
		otelzap.L().Ctx(ctx).Error("download avatar failed")
		return nil, fmt.Errorf("download avatar failed, status code: %d", resp.StatusCode)
	}

	if resp.IsSuccessState() {
		return resp.Body, nil
	}

	return nil, fmt.Errorf("download avatar failed, status code: %d", resp.StatusCode)
}
