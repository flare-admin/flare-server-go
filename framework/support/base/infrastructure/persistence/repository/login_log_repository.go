package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/mapper"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/entity"
)

type ILoginLogRepo interface {
	Create(ctx context.Context, log *entity.LoginLog) error
	// 动态查询方法
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*entity.LoginLog, error)
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}

type loginLogRepository struct {
	db     ILoginLogRepo
	mapper *mapper.LoginLogMapper
}

func NewLoginLogRepository(db ILoginLogRepo) repository.ILoginLogRepository {
	return &loginLogRepository{
		db:     db,
		mapper: &mapper.LoginLogMapper{},
	}
}

func (r *loginLogRepository) Create(ctx context.Context, log *model.LoginLog) error {
	toEntity := r.mapper.ToEntity(log)
	return r.db.Create(ctx, toEntity)
}
