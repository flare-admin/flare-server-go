package baserepo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
)

// SupportedIDTypes 支持的 ID 类型
type SupportedIDTypes interface {
	constraints.Integer | ~string
}

// IModel 数据模型接口
type IModel interface {
	GetPrimaryKey() string
	TableName() string
}

// IBaseRepo 基础仓库接口
type IBaseRepo[T IModel, I SupportedIDTypes] interface {
	FindById(ctx context.Context, id I) (*T, error)
	FindByIds(ctx context.Context, ids []I) ([]*T, error)
	DelById(ctx context.Context, id I) error
	DelByIds(ctx context.Context, ids []I) error
	DelByIdUnScoped(ctx context.Context, id I) error
	DelByIdsUnScoped(ctx context.Context, ids []I) error
	EditById(ctx context.Context, id I, data *T) error
	Add(ctx context.Context, data *T) (*T, error)
	BathAdd(ctx context.Context, data ...*T) error
	Count(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
	Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*T, error)
	Db(ctx context.Context) *gorm.DB
	GetDb() database.IDataBase
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	InIndependentTx(ctx context.Context, fn func(ctx context.Context) error) error
	GenStringId() string
	GenInt64Id() int64
}

// BaseRepo 基础仓库实现
type BaseRepo[T IModel, I SupportedIDTypes] struct {
	Model T
	db    database.IDataBase
}

func NewBaseRepo[T IModel, I SupportedIDTypes](db database.IDataBase, model T) *BaseRepo[T, I] {
	return &BaseRepo[T, I]{db: db, Model: model}
}

// --------------------------- 基础查询 ---------------------------

func (r *BaseRepo[T, I]) FindById(ctx context.Context, id I) (*T, error) {
	var res T
	db := r.db.DB(ctx).Model(r.Model)
	if hasDeletedField(r.Model) {
		db = db.Where("deleted_at = 0")
	}
	if err := db.Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *BaseRepo[T, I]) FindByIds(ctx context.Context, ids []I) ([]*T, error) {
	var res []*T
	db := r.db.DB(ctx).Model(r.Model)
	if hasDeletedField(r.Model) {
		db = db.Where("deleted_at = 0")
	}
	if err := db.Where(fmt.Sprintf("%s in ?", r.Model.GetPrimaryKey()), ids).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// --------------------------- 删除 ---------------------------

func (r *BaseRepo[T, I]) DelById(ctx context.Context, id I) error {
	db := r.db.DB(ctx).Model(r.Model)
	if hasDeletedField(r.Model) {
		if r.Model.GetPrimaryKey() == "" {
			return errors.New("primary key not defined")
		}
		return db.Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id).Update("deleted_at", utils.GetDateUnixMilli()).Error
	}
	return db.Delete(r.Model, id).Error
}

func (r *BaseRepo[T, I]) DelByIds(ctx context.Context, ids []I) error {
	db := r.db.DB(ctx).Model(r.Model)
	if hasDeletedField(r.Model) {
		if r.Model.GetPrimaryKey() == "" {
			return errors.New("primary key not defined")
		}
		return db.Where(fmt.Sprintf("%s in ?", r.Model.GetPrimaryKey()), ids).Update("deleted_at", utils.GetDateUnixMilli()).Error
	}
	return db.Where(fmt.Sprintf("%s in ?", r.Model.GetPrimaryKey()), ids).Delete(r.Model).Error
}

func (r *BaseRepo[T, I]) DelByIdUnScoped(ctx context.Context, id I) error {
	return r.db.DB(ctx).Unscoped().Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id).Delete(r.Model).Error
}

func (r *BaseRepo[T, I]) DelByIdsUnScoped(ctx context.Context, ids []I) error {
	return r.db.DB(ctx).Unscoped().Where(fmt.Sprintf("%s in ?", r.Model.GetPrimaryKey()), ids).Delete(r.Model).Error
}

// --------------------------- 更新 ---------------------------

func (r *BaseRepo[T, I]) EditById(ctx context.Context, id I, data *T) error {
	if data == nil {
		return errors.New("update data is nil")
	}
	db := r.db.DB(ctx).Model(r.Model)
	if hasDeletedField(r.Model) {
		db = db.Where("deleted_at = 0")
	}
	db = db.Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id)

	setUpdatedAt(data)
	return db.Updates(data).Error
}

// --------------------------- 添加 ---------------------------

func (r *BaseRepo[T, I]) Add(ctx context.Context, data *T) (*T, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	setCreatedAt(data)
	setIDIfEmpty(r, data)
	if err := r.db.DB(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *BaseRepo[T, I]) BathAdd(ctx context.Context, data ...*T) error {
	if len(data) == 0 {
		return nil
	}
	for _, d := range data {
		setCreatedAt(d)
		setIDIfEmpty(r, d)
	}
	return r.db.DB(ctx).Create(data).Error
}

// --------------------------- 查询 ---------------------------

func (r *BaseRepo[T, I]) Count(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	var count int64
	db := r.db.DB(ctx).Model(r.Model)
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}
	return count, db.Count(&count).Error
}

func (r *BaseRepo[T, I]) Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*T, error) {
	var res []*T
	db := r.db.DB(ctx).Model(r.Model)
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}
	if order := qb.BuildOrderBy(); order != "" {
		db = db.Order(order)
	}
	if limit, vals := qb.BuildLimit(); limit != "" {
		db = db.Offset(vals[0]).Limit(vals[1])
	}
	err := db.Find(&res).Error
	if database.IfErrorNotFound(err) {
		return res, nil
	}
	return res, err
}

// --------------------------- 事务 & DB ---------------------------

func (r *BaseRepo[T, I]) Db(ctx context.Context) *gorm.DB {
	return r.db.DB(ctx)
}

func (r *BaseRepo[T, I]) GetDb() database.IDataBase {
	return r.db
}

func (r *BaseRepo[T, I]) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.InTx(ctx, fn)
}

func (r *BaseRepo[T, I]) InIndependentTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.InIndependentTx(ctx, fn)
}

func (r *BaseRepo[T, I]) GenStringId() string {
	return r.db.GenStringId()
}

func (r *BaseRepo[T, I]) GenInt64Id() int64 {
	return r.db.GenInt64Id()
}

// --------------------------- 工具函数 ---------------------------

// 检查是否有 DeletedAt 字段
func hasDeletedField[T IModel](model T) bool {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	_, ok := v.Type().FieldByName("DeletedAt")
	return ok
}

// 设置 UpdatedAt
func setUpdatedAt[T IModel](data *T) {
	v := reflect.ValueOf(data).Elem()
	if f := v.FieldByName("UpdatedAt"); f.IsValid() && f.CanSet() {
		f.SetInt(time.Now().UnixMilli())
	}
}

// 设置 CreatedAt
func setCreatedAt[T IModel](data *T) {
	v := reflect.ValueOf(data).Elem()
	if f := v.FieldByName("CreatedAt"); f.IsValid() && f.CanSet() {
		f.SetInt(time.Now().UnixMilli())
	}
}

// 自动生成 ID
func setIDIfEmpty[T IModel, I SupportedIDTypes](r *BaseRepo[T, I], data *T) {
	v := reflect.ValueOf(data).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		if field.Name == "ID" || field.Name == "Id" || field.Tag.Get("pk") == "true" {
			isEmpty := false
			switch f.Kind() {
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32:
				isEmpty = f.Int() == 0
				if isEmpty {
					f.SetInt(r.db.GenInt64Id())
				}
			case reflect.String:
				isEmpty = f.String() == ""
				if isEmpty {
					f.SetString(r.db.GenStringId())
				}
			}
			break
		}
	}
}
