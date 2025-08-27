package oplog

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/repository"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
)

type DbOperationLogWriter struct {
	repo repository.IOperationLogRepo
}

func NewDbOperationLogWriter(repo repository.IOperationLogRepo) oplog.IDbOperationLogWrite {
	return &DbOperationLogWriter{
		repo: repo,
	}
}

func (w *DbOperationLogWriter) Save(ctx context.Context, data *oplog.OperationLog) error {
	log := &entity.OperationLog{
		UserID:    data.UserID,
		Username:  data.Username,
		TenantID:  data.TenantID,
		Method:    data.Method,
		Path:      data.Path,
		Query:     data.Query,
		Body:      data.Body,
		IP:        data.IP,
		UserAgent: data.UserAgent,
		Status:    data.Status,
		Error:     data.Error,
		Duration:  data.Duration,
		Module:    data.Module,
		Action:    data.Action,
		BaseIntTime: database.BaseIntTime{
			CreatedAt: data.CreatedAt.Unix(),
		},
	}
	return w.repo.Create(ctx, log)
}
