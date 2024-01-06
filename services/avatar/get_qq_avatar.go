package avatar

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kolesa-team/go-webp/webp"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/cloudwego/hertz/pkg/common/bytebufferpool"
	"github.com/nfnt/resize"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func (s *service) GetQQAvatar(ctx context.Context, hash string, args GetAvatarArgs) (io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.GetQQAvatar")
	defer span.End()

	// get qq avatar
	qqid, err := s.MD5QQMappingRepo.GetQQIdByEmailMD5(ctx, hash)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("get qq id by email md5 failed", zap.Error(err))
		return nil, err
	}

	stream, err := s.downloadAvatar(ctx, fmt.Sprintf("https://q.qlogo.cn/headimg_dl?dst_uin=%d&spec=640&img_type=jpg", qqid))
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

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if args.EnableWebp {
		if err := webp.Encode(buf, img, nil); err != nil {
			otelzap.L().Ctx(ctx).Error("encode qq avatar failed", zap.Error(err))
			return nil, err
		}
	} else {
		if err := png.Encode(buf, img); err != nil {
			otelzap.L().Ctx(ctx).Error("encode qq avatar failed", zap.Error(err))
			return nil, err
		}
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
