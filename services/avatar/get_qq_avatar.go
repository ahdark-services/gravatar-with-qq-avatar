package avatar

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"strconv"

	"github.com/nfnt/resize"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func initQQClient() *req.Client {
	c := req.C().
		EnableDumpEachRequest().
		SetBaseURL("https://q.qlogo.cn/").
		SetCommonHeader("Accept", "image/jpeg").
		SetCommonQueryParams(map[string]string{
			"dst_uin":  "0",
			"spec":     "640",
			"img_type": "jpg",
		}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error (SetSuccessResult or SetErrorResult was invoked before).
				if dump := resp.Dump(); dump != "" { // Append dump content to original underlying error to help troubleshoot.
					resp.Err = fmt.Errorf("error: %s\nraw content:\n%s", resp.Err.Error(), resp.Dump())
				}
				return nil // Skip the following logic if there is an underlying error.
			}

			if !resp.IsSuccessState() {
				resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
			}

			return nil
		}).
		WrapRoundTripFunc(WithTracer)

	return c
}

func (s *service) getQQAvatar(ctx context.Context, hash string, args GetAvatarArgs) (image.Image, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.getQQAvatar")
	defer span.End()

	// get qq avatar
	qqid, err := s.MD5QQMappingRepo.GetQQIdByEmailMD5(ctx, hash)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("get qq id by email md5 failed", zap.Error(err))
		return nil, err
	}

	resp, err := s.qqAvatarClient.R().
		SetContext(ctx).
		SetQueryParam("dst_uin", strconv.FormatInt(qqid, 10)).
		Get("headimg_dl")
	if err != nil {
		otelzap.L().Ctx(ctx).Error("download qq avatar failed", zap.Error(err))
		return nil, err
	} else if resp.IsErrorState() {
		otelzap.L().Ctx(ctx).Error("download qq avatar failed", zap.Error(err))
		return nil, fmt.Errorf("download qq avatar failed: %v", resp.ErrorResult())
	}

	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("decode qq avatar failed", zap.Error(err))
		return nil, err
	}

	img = resize.Resize(uint(args.Size), uint(args.Size), img, resize.Lanczos3)

	return img, nil
}
