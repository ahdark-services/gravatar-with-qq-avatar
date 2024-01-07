package avatar

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"image"
	"strconv"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type GetGravatarResult struct {
	Avatar       image.Image
	LastModified time.Time
}

func initGravatarClient() *req.Client {
	c := req.C().
		EnableDumpEachRequest().
		SetBaseURL("https://gravatar.com/avatar/").
		SetCommonHeader("Accept", "image/jpeg").
		SetCommonQueryParams(map[string]string{
			"s": "80",
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

func (s *service) getGravatar(ctx context.Context, hash string, args GetAvatarArgs) (GetGravatarResult, error) {
	ctx, span := tracer.Start(ctx, "service.AvatarService.getGravatar")
	defer span.End()

	resp, err := s.gravatarClient.R().
		SetContext(ctx).
		SetPathParam("hash", lo.If(args.ForceDefault, "").Else(hash)).
		SetQueryParams(map[string]string{
			"d": args.Default,
			"s": strconv.FormatInt(args.Size, 10),
			"r": args.Rating,
		}).
		Get("{hash}")
	if resp.GetStatusCode() == 404 {
		return GetGravatarResult{}, nil
	} else if err != nil {
		otelzap.L().Ctx(ctx).Error("download gravatar failed", zap.Error(err))
		return GetGravatarResult{}, err
	} else if resp.IsErrorState() {
		otelzap.L().Ctx(ctx).Error("download gravatar failed", zap.Error(resp.Err))
		return GetGravatarResult{}, fmt.Errorf("download gravatar failed: %v", resp.ErrorResult())
	}

	defer resp.Body.Close()

	img, err := parseImage(resp.GetHeader("Content-Type"), resp.Body)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("parse image failed", zap.Error(err))
		return GetGravatarResult{}, err
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
