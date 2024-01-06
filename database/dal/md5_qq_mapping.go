package dal

import (
	"context"
	"github.com/scylladb/gocqlx/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AH-dark/gravatar-with-qq-support/database/models"
)

type MD5QQMappingRepo interface {
	GetQQIdByEmailMD5(ctx context.Context, emailMD5 string) (int64, error)
	InsertMapping(ctx context.Context, qqid int64, emailMD5 string) error
}

type MD5QQMappingRepoImpl struct {
	fx.In
	Session *gocqlx.Session
}

func NewMD5QQMapping(repo MD5QQMappingRepoImpl) MD5QQMappingRepo {
	return &repo
}

func (repo *MD5QQMappingRepoImpl) GetQQIdByEmailMD5(ctx context.Context, emailMD5 string) (int64, error) {
	ctx, span := tracer.Start(ctx, "dal.MD5QQMappingRepo.GetQQIdByEmailMD5")
	defer span.End()

	var qqId int64
	if err := models.MD5QQMappingTable.
		GetQueryContext(ctx, *repo.Session, "qq_id").
		BindStruct(&models.MD5QQMapping{EmailMD5: emailMD5}).
		Scan(&qqId); err != nil {
		otelzap.L().Ctx(ctx).Error("get qq id by email md5 failed", zap.Error(err))
		return 0, err
	}

	return qqId, nil
}

func (repo *MD5QQMappingRepoImpl) InsertMapping(ctx context.Context, qqid int64, emailMD5 string) error {
	ctx, span := tracer.Start(ctx, "dal.MD5QQMappingRepo.InsertMapping")
	defer span.End()

	if err := models.MD5QQMappingTable.
		InsertQueryContext(ctx, *repo.Session).
		BindStruct(&models.MD5QQMapping{
			EmailMD5: emailMD5,
			QQId:     qqid,
		}).
		ExecRelease(); err != nil {
		otelzap.L().Ctx(ctx).Error("insert qq avatar failed", zap.Error(err))
		return err
	}

	return nil
}
