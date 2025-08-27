package data

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/repository"
)

// categoryRepository 分类数据访问实现
type categoryRepository struct {
	*baserepo.BaseRepo[entity.Category, string]
}

// NewCategoryRepository 创建分类数据访问实例
func NewCategoryRepository(data database.IDataBase) repository.ICategoryRepository {
	// 同步表
	tables := []interface{}{
		&entity.Category{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables error: %v", err)
	}
	return &categoryRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.Category, string](data, entity.Category{}),
	}
}

func (c *categoryRepository) FindByCode(ctx context.Context, code string) (*entity.Category, error) {
	var category entity.Category
	err := c.Db(ctx).Model(&entity.Category{}).Where("code = ?", code).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *categoryRepository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	var categories []*entity.Category
	err := c.Db(ctx).Model(&entity.Category{}).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
