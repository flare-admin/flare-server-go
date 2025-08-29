package data

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/repository"
)

// templateRepository 模板仓储实现
type templateRepository struct {
	*baserepo.BaseRepo[entity.Template, string]
}

// NewTemplateRepository 创建模板仓储
func NewTemplateRepository(data database.IDataBase) repository.ITemplateRepository {
	// 同步表
	tables := []interface{}{
		&entity.Template{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables error: %v", err)
	}
	return &templateRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.Template, string](data),
	}
}

// FindByCode 根据编码查询模板
func (t *templateRepository) FindByCode(ctx context.Context, code string) (*entity.Template, error) {
	var template entity.Template
	err := t.Db(ctx).Model(&entity.Template{}).Where("code = ?", code).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByCategoryID 根据分类ID查询模板列表
func (t *templateRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.Template, error) {
	template := make([]*entity.Template, 0)
	err := t.Db(ctx).Model(&entity.Template{}).Where("category_id = ?", categoryID).Find(&template).Error
	if err != nil {
		return nil, err
	}
	return template, nil
}

// ExistsByCode 检查编码是否存在
func (t *templateRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := t.Db(ctx).Model(&entity.Template{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查询所有模板
func (t *templateRepository) FindAll(ctx context.Context) ([]*entity.Template, error) {
	templates := make([]*entity.Template, 0)
	err := t.Db(ctx).Model(&entity.Template{}).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
