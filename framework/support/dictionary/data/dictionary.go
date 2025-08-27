package data

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/model"

	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type IDictionaryRepo interface {
	baserepo.IBaseRepo[model.Category, string]
	ListCategories(ctx context.Context) ([]*model.Category, error)
	CreateOption(ctx context.Context, option *model.Option) error
	UpdateOption(ctx context.Context, option *model.Option) error
	DeleteOption(ctx context.Context, id string) error
	GetOptions(ctx context.Context, categoryID, Keyword string, status *int) ([]*model.Option, error)
	FindOptionById(ctx context.Context, id string) (*model.Option, error)
}

// categoryRepo ， 菜单数据层
type categoryRepo struct {
	*baserepo.BaseRepo[model.Category, string]
}

// NewDictionaryRepo ，
// 参数：
//
//	data ： desc
//
// 返回值：
//
//	biz.ISysMenuRepo ：desc
func NewDictionaryRepo(data database.IDataBase) IDictionaryRepo {
	// 同步表
	tables := []interface{}{
		&model.Category{},
		&model.Option{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &categoryRepo{
		BaseRepo: baserepo.NewBaseRepo[model.Category, string](data, model.Category{}),
	}
}

// ListCategories 列出所有分类
func (r *categoryRepo) ListCategories(ctx context.Context) ([]*model.Category, error) {
	var res []*model.Category
	err := r.Db(ctx).Find(&res).Error
	return res, err
}

// CreateOption 创建选项
func (r *categoryRepo) CreateOption(ctx context.Context, option *model.Option) error {
	return r.Db(ctx).Create(option).Error
}

// UpdateOption 更新选项
func (r *categoryRepo) UpdateOption(ctx context.Context, option *model.Option) error {
	err := r.Db(ctx).Save(option).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteOption 删除选项
func (r *categoryRepo) DeleteOption(ctx context.Context, id string) error {
	err := r.Db(ctx).Unscoped().Delete(&model.Option{ID: id}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetOptions 获取选项列表
func (r *categoryRepo) GetOptions(ctx context.Context, categoryID, Keyword string, status *int) ([]*model.Option, error) {
	var res []*model.Option
	err := r.Db(ctx).Scopes(func(d *gorm.DB) *gorm.DB {
		if categoryID != "" {
			d = d.Where("category_id = ?", categoryID)
		}
		if Keyword != "" {
			d = d.Where("code like ? or value like ? or label like ?", "%"+Keyword+"%", "%"+Keyword+"%", "%"+Keyword+"%")
		}
		if status != nil {
			d = d.Where("status = ?", status)
		}
		return d
	}).Order("sort").Find(&res).Error
	return res, err
}

// FindOptionById 根据id查询选项
func (r *categoryRepo) FindOptionById(ctx context.Context, id string) (*model.Option, error) {
	var res model.Option
	err := r.Db(ctx).Where("id = ?", id).First(&res).Error
	return &res, err
}
