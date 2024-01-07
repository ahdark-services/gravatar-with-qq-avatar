package avatar

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func (s *service) getQQAvatar(ctx context.Context, hash string, args GetAvatarArgs) (image.Image, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.getQQAvatar")
	defer span.End()

	// get qq avatar
	qqid, err := s.MD5QQMappingRepo.GetQQIdByEmailMD5(ctx, hash)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("get qq id by email md5 failed", zap.Error(err))
		return nil, err
	}

	stream, _, err := s.downloadAvatar(ctx, fmt.Sprintf("https://q.qlogo.cn/headimg_dl?dst_uin=%d&spec=640&img_type=jpg", qqid))
	if err != nil {
		otelzap.L().Ctx(ctx).Error("download qq avatar failed", zap.Error(err))
		return nil, err
	}

	img, err := jpeg.Decode(stream)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("decode qq avatar failed", zap.Error(err))
		return nil, err
	}

	if err := stream.Close(); err != nil {
		otelzap.L().Ctx(ctx).Warn("close stream failed", zap.Error(err))
	}

	img = resize.Resize(uint(args.Size), uint(args.Size), img, resize.Lanczos3)

	return img, nil
}
