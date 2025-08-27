package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/entity"
)

type IOperationLogRepo interface {
	Create(ctx context.Context, log *entity.OperationLog) error
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*entity.OperationLog, error)
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}
type operationLogRepository struct {
	db database.IDataBase
}

func NewOperationLogRepository(db database.IDataBase) IOperationLogRepo {
	return &operationLogRepository{
		db: db,
	}
}

// Create 创建操作日志
func (r *operationLogRepository) Create(ctx context.Context, log *entity.OperationLog) error {
	t := time.Unix(log.CreatedAt, 0)
	// 确保表存在
	if err := r.EnsureTable(ctx, log.TenantID, t); err != nil {
		return err
	}

	return r.db.DB(ctx).Table(r.GetTableName(log.TenantID, t)).Create(log).Error
}

// Find 查询操作日志列表
func (r *operationLogRepository) Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*entity.OperationLog, error) {
	// 确保表存在
	if err := r.EnsureTable(ctx, tenantID, month); err != nil {
		return nil, err
	}

	var entities []*entity.OperationLog
	db := r.db.DB(ctx).Table(r.GetTableName(tenantID, month))

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	// 添加排序
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	}

	// 添加分页
	if limit, offset := qb.BuildLimit(); limit != "" {
		db = db.Limit(offset[1]).Offset(offset[0])
	}

	// 执行查询
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// Count 统计数量
func (r *operationLogRepository) Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error) {
	// 确保表存在
	if err := r.EnsureTable(ctx, tenantID, month); err != nil {
		return 0, err
	}

	var count int64
	db := r.db.DB(ctx).Table(r.GetTableName(tenantID, month))

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// EnsureTable 确保表存在
func (r *operationLogRepository) EnsureTable(ctx context.Context, tenantID string, month time.Time) error {
	tableName := r.GetTableName(tenantID, month)

	// 检查表是否存在
	if r.db.DB(ctx).Migrator().HasTable(tableName) {
		return nil
	}

	// 创建表结构
	type OperationLogTable struct {
		entity.OperationLog
	}

	// 使用 GORM 自动迁移创建表
	return r.db.DB(ctx).Table(tableName).AutoMigrate(&OperationLogTable{})
}

// GetTableName 获取表名
func (r *operationLogRepository) GetTableName(tenantID string, month time.Time) string {
	return fmt.Sprintf("sys_operation_log_%s_%s", tenantID, month.Format("200601"))
}
