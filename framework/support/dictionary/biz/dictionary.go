package biz

import (
	"context"
	"errors"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/data"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/dto"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/model"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/service"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/translator"
)

type DictionaryUseCase struct {
	repo  data.IDictionaryRepo
	tr    translator.ITranslator
	idGen snowflake_id.IIdGenerate
}

func NewDictionaryUseCase(repo data.IDictionaryRepo, tr translator.ITranslator, idGen snowflake_id.IIdGenerate) service.IDictionaryService {
	return &DictionaryUseCase{repo: repo, tr: tr, idGen: idGen}
}

// CreateCategory 创建分类
func (uc *DictionaryUseCase) CreateCategory(ctx context.Context, req *dto.CategoryCreateReq) herrors.Herr {
	// 业务校验
	if req.ID == "" || req.Name == "" {
		return herrors.CreateFail(errors.New("分类ID和名称不能为空"))
	}

	// 检查ID是否已存在
	if _, err := uc.repo.FindById(ctx, req.ID); err == nil {
		return herrors.CreateFail(errors.New("分类ID已存在"))
	}

	category := &model.Category{
		ID:          req.ID,
		Name:        req.Name,
		I18nKey:     req.I18nKey,
		Description: req.Description,
	}
	_, err := uc.repo.Add(ctx, category)
	if err != nil {
		return herrors.CreateFail(err)
	}
	return nil
}

// UpdateCategory 更新分类
func (uc *DictionaryUseCase) UpdateCategory(ctx context.Context, id string, req *dto.CategoryUpdateReq) herrors.Herr {
	// 检查分类是否存在
	if _, err := uc.repo.FindById(ctx, id); err != nil {
		return herrors.UpdateFail(errors.New("分类不存在"))
	}

	category := &model.Category{
		Name:        req.Name,
		I18nKey:     req.I18nKey,
		Description: req.Description,
	}
	err := uc.repo.EditById(ctx, category)
	if err != nil {
		return herrors.UpdateFail(err)
	}
	return nil
}

// DelCategory 删除分类
func (uc *DictionaryUseCase) DelCategory(ctx context.Context, id string) herrors.Herr {
	if id == "" {
		return herrors.DeleteFail(errors.New("分类ID不能为空"))
	}
	err := uc.repo.DelByIdUnScoped(ctx, id)
	if err != nil {
		return herrors.DeleteFail(err)
	}
	// 清除缓存
	uc.tr.ClearCache(ctx, id)
	return nil
}

// GetCategory 获取分类
func (uc *DictionaryUseCase) GetCategory(ctx context.Context, id string) (*dto.Category, herrors.Herr) {
	if id == "" {
		return nil, herrors.QueryFail(errors.New("分类ID不能为空"))
	}
	category, err := uc.repo.FindById(ctx, id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return modelToDTO(category), nil
}

// ListCategories 列出所有分类
func (uc *DictionaryUseCase) ListCategories(ctx context.Context, req *dto.CategoryQueryReq) ([]*dto.Category, int64, herrors.Herr) {
	query := db_query.NewQueryBuilder()
	if req.Id != "" {
		query.Where("id", db_query.Eq, req.Id)
	}
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.I18nKey != "" {
		query.Where("i18n_key", db_query.Eq, req.I18nKey)
	}

	query.WithPage(&req.Page)

	total, err := uc.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}
	cars, err := uc.repo.Find(ctx, query)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}
	res := make([]*dto.Category, 0, len(cars))
	for _, v := range cars {
		res = append(res, modelToDTO(v))
	}
	return res, total, nil
}

// GetOptions 获取选项列表
func (uc *DictionaryUseCase) GetOptions(ctx context.Context, req *dto.OptionQueryReq) ([]*dto.Option, herrors.Herr) {
	if req.CategoryID != "" {
		// 验证分类是否存在
		if _, err := uc.repo.FindById(ctx, req.CategoryID); err != nil {
			return nil, herrors.QueryFail(errors.New("分类不存在"))
		}
	}

	options, err := uc.repo.GetOptions(ctx, req.CategoryID, req.Keyword, req.Status)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换结果
	dtoOptions := make([]*dto.Option, 0, len(options))
	for _, opt := range options {
		if dtoOpt := optionModelToDTO(opt); dtoOpt != nil {
			dtoOptions = append(dtoOptions, dtoOpt)
		}
	}
	return dtoOptions, nil
}

// CreateOption 创建选项
func (uc *DictionaryUseCase) CreateOption(ctx context.Context, req *dto.OptionCreateReq) herrors.Herr {
	// 业务校验
	if req.CategoryID == "" || req.Value == "" {
		return herrors.CreateFail(errors.New("选项ID、分类ID、编码和值不能为空"))
	}

	// 验证分类是否存在
	if _, err := uc.repo.FindById(ctx, req.CategoryID); err != nil {
		return herrors.CreateFail(errors.New("所属分类不存在"))
	}

	// 检查选项ID是否已存在
	options, err := uc.repo.GetOptions(ctx, req.CategoryID, "", nil)
	if err != nil {
		return herrors.CreateFail(err)
	}

	for _, opt := range options {
		if opt.Value == req.Value {
			return herrors.CreateFail(errors.New("选项已存在"))
		}
	}

	option := &model.Option{
		ID:         uc.idGen.GenStringId(),
		CategoryID: req.CategoryID,
		Value:      req.Value,
		I18nKey:    req.I18nKey,
		Sort:       req.Sort,
		Status:     req.Status,
		Remark:     req.Remark,
		Label:      req.Label,
	}
	err = uc.repo.CreateOption(ctx, option)
	if err != nil {
		return herrors.CreateFail(err)
	}
	uc.tr.ClearCache(ctx, req.CategoryID)
	return nil
}

// UpdateOption 更新选项
func (uc *DictionaryUseCase) UpdateOption(ctx context.Context, id string, req *dto.OptionUpdateReq) herrors.Herr {
	options, err := uc.repo.GetOptions(ctx, "", "", nil)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	var found *model.Option
	for _, opt := range options {
		if opt.ID == id {
			found = opt
			break
		}
	}
	if found == nil {
		return herrors.UpdateFail(errors.New("选项不存在"))
	}

	// 检查编码是否重复
	if req.Value != "" && req.Value != found.Value {
		for _, opt := range options {
			if opt.ID != id && opt.CategoryID == found.CategoryID && opt.Value == req.Value {
				return herrors.UpdateFail(errors.New("选项已存在"))
			}
		}
	}

	if req.Value != "" {
		found.Value = req.Value
	}
	if req.I18nKey != "" {
		found.I18nKey = req.I18nKey
	}
	if req.Label != "" {
		found.Label = req.Label
	}
	found.Sort = req.Sort
	found.Status = req.Status
	if req.Remark != "" {
		found.Remark = req.Remark
	}

	err = uc.repo.UpdateOption(ctx, found)
	if err != nil {
		return herrors.UpdateFail(err)
	}
	uc.tr.ClearCache(ctx, found.CategoryID)
	return nil
}

// DeleteOption 删除选项
func (uc *DictionaryUseCase) DeleteOption(ctx context.Context, id string) herrors.Herr {
	if id == "" {
		return herrors.DeleteFail(errors.New("选项ID不能为空"))
	}
	found, err := uc.repo.FindOptionById(ctx, id)
	if err != nil {
		return herrors.DeleteFail(err)
	}
	err = uc.repo.DeleteOption(ctx, id)
	if err != nil {
		return herrors.DeleteFail(err)
	}
	uc.tr.ClearCache(ctx, found.CategoryID)
	return nil
}

// modelToDTO 转换方法
func modelToDTO(m *model.Category) *dto.Category {
	if m == nil {
		return nil
	}
	return &dto.Category{
		BaseModel: database.BaseModel{
			BaseIntTime: m.BaseIntTime,
			Creator:     m.Creator,
			Updater:     m.Updater,
		},
		ID:          m.ID,
		Name:        m.Name,
		I18nKey:     m.I18nKey,
		Description: m.Description,
	}
}

func optionModelToDTO(m *model.Option) *dto.Option {
	if m == nil {
		return nil
	}
	return &dto.Option{
		ID:         m.ID,
		CategoryID: m.CategoryID,
		Value:      m.Value,
		I18nKey:    m.I18nKey,
		Sort:       m.Sort,
		Status:     m.Status,
		Remark:     m.Remark,
		Label:      m.Label,
		BaseModel: database.BaseModel{
			BaseIntTime: m.BaseIntTime,
			Creator:     m.Creator,
			Updater:     m.Updater,
		},
	}
}
