package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/dto"
)

// IDictionaryService 字典服务接口
type IDictionaryService interface {
	// Category 分类相关接口
	CreateCategory(ctx context.Context, req *dto.CategoryCreateReq) herrors.Herr
	UpdateCategory(ctx context.Context, id string, req *dto.CategoryUpdateReq) herrors.Herr
	GetCategory(ctx context.Context, id string) (*dto.Category, herrors.Herr)
	DelCategory(ctx context.Context, id string) herrors.Herr
	ListCategories(ctx context.Context, req *dto.CategoryQueryReq) ([]*dto.Category, int64, herrors.Herr)

	// Option 选项相关接口
	CreateOption(ctx context.Context, req *dto.OptionCreateReq) herrors.Herr
	UpdateOption(ctx context.Context, id string, req *dto.OptionUpdateReq) herrors.Herr
	DeleteOption(ctx context.Context, id string) herrors.Herr
	GetOptions(ctx context.Context, req *dto.OptionQueryReq) ([]*dto.Option, herrors.Herr)
}
