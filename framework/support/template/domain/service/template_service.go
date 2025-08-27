package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/constants"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/valueobject"
	"regexp"
	"time"
)

// TemplateService 模板领域服务
type TemplateService struct {
	templateRepo repository.TemplateRepository
	categoryRepo repository.CategoryRepository
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(
	templateRepo repository.TemplateRepository,
	categoryRepo repository.CategoryRepository,
) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(ctx context.Context, cmd *valueobject.CreateTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 验证分类是否存在
	category, err := s.categoryRepo.FindByID(ctx, cmd.CategoryID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if category == nil {
		return template_err.CategoryNotExist
	}

	// 3. 验证模板编码是否已存在
	exists, err := s.templateRepo.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "检查模板编码是否存在失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if exists {
		return template_err.TemplateCodeExists
	}

	// 4. 创建模板
	template := model.NewTemplate(cmd.Code, cmd.Name, cmd.Description, cmd.CategoryID)
	for _, attr := range cmd.Attributes {
		template.AddAttribute(convertToModelAttribute(attr))
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		hlog.CtxErrorf(ctx, "创建模板失败:%v", err)
		return template_err.TemplateCreateFailed(err)
	}

	hlog.CtxInfof(ctx, "创建模板成功,编码:%s,名称:%s", template.Code, template.Name)
	return nil
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, cmd *valueobject.UpdateTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取原模板信息
	oldTemplate, err := s.templateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if oldTemplate == nil {
		return template_err.TemplateNotExist
	}

	// 3. 验证模板状态
	if oldTemplate.Status != 1 {
		hlog.CtxErrorf(ctx, "模板已禁用,状态:%d", oldTemplate.Status)
		return template_err.TemplateDisabled
	}

	// 4. 验证模板编码是否已存在（如果编码发生变化）
	if cmd.Code != oldTemplate.Code {
		exists, err := s.templateRepo.ExistsByCode(ctx, cmd.Code)
		if err != nil {
			hlog.CtxErrorf(ctx, "检查模板编码是否存在失败:%v", err)
			return template_err.TemplateGetFailed(err)
		}
		if exists {
			return template_err.TemplateCodeExists
		}
	}

	// 5. 验证分类是否存在
	if cmd.CategoryID != oldTemplate.CategoryID {
		category, err := s.categoryRepo.FindByID(ctx, cmd.CategoryID)
		if err != nil {
			hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
			return template_err.TemplateGetFailed(err)
		}
		if category == nil {
			return template_err.CategoryNotExist
		}
	}

	// 6. 更新模板
	oldTemplate.Code = cmd.Code
	oldTemplate.Name = cmd.Name
	oldTemplate.Description = cmd.Description
	oldTemplate.CategoryID = cmd.CategoryID
	oldTemplate.Attributes = make([]model.Attribute, len(cmd.Attributes))
	for i, attr := range cmd.Attributes {
		oldTemplate.Attributes[i] = convertToModelAttribute(attr)
	}
	oldTemplate.UpdatedAt = utils.GetDateUnix()

	if err := s.templateRepo.Update(ctx, oldTemplate); err != nil {
		hlog.CtxErrorf(ctx, "更新模板失败:%v", err)
		return template_err.TemplateUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "更新模板成功,模板ID:%s,名称:%s", oldTemplate.ID, oldTemplate.Name)
	return nil
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(ctx context.Context, cmd *valueobject.DeleteTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取模板信息
	template, err := s.templateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if template == nil {
		return template_err.TemplateNotExist
	}

	// 3. 验证模板状态
	if template.Status != 1 {
		hlog.CtxErrorf(ctx, "模板已禁用,状态:%d", template.Status)
		return template_err.TemplateDisabled
	}

	// 4. 删除模板
	if err := s.templateRepo.Delete(ctx, cmd.ID); err != nil {
		hlog.CtxErrorf(ctx, "删除模板失败:%v", err)
		return template_err.TemplateDeleteFailed(err)
	}

	hlog.CtxInfof(ctx, "删除模板成功,模板ID:%s", cmd.ID)
	return nil
}

// EnableTemplate 启用模板
func (s *TemplateService) EnableTemplate(ctx context.Context, cmd *valueobject.UpdateTemplateStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取模板信息
	template, err := s.templateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if template == nil {
		return template_err.TemplateNotExist
	}

	// 3. 更新状态
	template.Status = 1
	template.UpdatedAt = utils.GetDateUnix()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		hlog.CtxErrorf(ctx, "更新模板状态失败:%v", err)
		return template_err.TemplateUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "启用模板成功,模板ID:%s", cmd.ID)
	return nil
}

// DisableTemplate 禁用模板
func (s *TemplateService) DisableTemplate(ctx context.Context, cmd *valueobject.UpdateTemplateStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取模板信息
	template, err := s.templateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if template == nil {
		return template_err.TemplateNotExist
	}

	// 3. 更新状态
	template.Status = 2
	template.UpdatedAt = utils.GetDateUnix()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		hlog.CtxErrorf(ctx, "更新模板状态失败:%v", err)
		return template_err.TemplateUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "禁用模板成功,模板ID:%s", cmd.ID)
	return nil
}

// FindByID 根据ID查询模板
func (s *TemplateService) FindByID(ctx context.Context, id string) (*model.Template, *herrors.HError) {
	template, err := s.templateRepo.FindByID(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, template_err.TemplateNotExist
		}
		return nil, template_err.TemplateGetFailed(err)
	}
	return template, nil
}

// ValidateTemplateContent 验证模板数据
func (s *TemplateService) ValidateTemplateContent(ctx context.Context, templateID string, content map[string]interface{}) *herrors.HError {
	// 1. 获取模板信息
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板失败:%v", err)
		return template_err.TemplateGetFailed(err)
	}
	if template == nil {
		return template_err.TemplateNotExist
	}

	// 2. 验证模板状态
	if template.Status != 1 {
		hlog.CtxErrorf(ctx, "模板已禁用,状态:%d", template.Status)
		return template_err.TemplateDisabled
	}

	// 3. 验证必填字段
	for _, attr := range template.Attributes {
		value, exists := content[attr.Key]
		if attr.Required && (!exists || value == nil) {
			return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]不能为空", attr.Name))
		}
		if !exists {
			continue
		}

		// 4. 根据属性类型验证数据
		if err := s.validateAttributeValue(attr, value); err != nil {
			return err
		}
	}

	return nil
}

// validateAttributeValue 验证属性值
func (s *TemplateService) validateAttributeValue(attr model.Attribute, value interface{}) *herrors.HError {
	switch attr.Type {
	case constants.AttributeTypeString, constants.AttributeTypeTextarea:
		return s.validateStringValue(attr, value)
	case constants.AttributeTypeNumber, constants.AttributeTypeMoney:
		return s.validateNumberValue(attr, value)
	case constants.AttributeTypeInteger:
		return s.validateIntegerValue(attr, value)
	case constants.AttributeTypeDate:
		return s.validateDateValue(attr, value)
	case constants.AttributeTypeDateTime:
		return s.validateDateTimeValue(attr, value)
	case constants.AttributeTypeTime:
		return s.validateTimeValue(attr, value)
	case constants.AttributeTypeSelect, constants.AttributeTypeSwitch:
		return s.validateSelectValue(attr, value)
	case constants.AttributeTypeBoolean:
		return s.validateBooleanValue(attr, value)
	case constants.AttributeTypeWallet:
		return s.validateWalletValue(attr, value)
	default:
		return template_err.TemplateValidationFailed(fmt.Errorf("不支持的属性类型: %s", attr.Type))
	}
}

// validateStringValue 验证字符串值
func (s *TemplateService) validateStringValue(attr model.Attribute, value interface{}) *herrors.HError {
	strValue, ok := value.(string)
	if !ok {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是字符串类型", attr.Name))
	}

	// 验证长度
	if attr.Validation.Length != nil {
		if len(strValue) != *attr.Validation.Length {
			return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]长度必须为%d", attr.Name, *attr.Validation.Length))
		}
	}

	// 验证正则表达式
	if attr.Validation.Pattern != "" {
		matched, err := regexp.MatchString(attr.Validation.Pattern, strValue)
		if err != nil {
			return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]正则表达式验证失败: %v", attr.Name, err))
		}
		if !matched {
			return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]格式不正确", attr.Name))
		}
	}

	return nil
}

// validateNumberValue 验证数字值
func (s *TemplateService) validateNumberValue(attr model.Attribute, value interface{}) *herrors.HError {
	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case int:
		numValue = float64(v)
	case int64:
		numValue = float64(v)
	default:
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是数字类型", attr.Name))
	}

	// 验证最小值
	if attr.Validation.Min != nil && numValue < *attr.Validation.Min {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]不能小于%v", attr.Name, *attr.Validation.Min))
	}

	// 验证最大值
	if attr.Validation.Max != nil && numValue > *attr.Validation.Max {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]不能大于%v", attr.Name, *attr.Validation.Max))
	}

	return nil
}

// validateIntegerValue 验证整数值
func (s *TemplateService) validateIntegerValue(attr model.Attribute, value interface{}) *herrors.HError {
	var intValue int64
	switch v := value.(type) {
	case int:
		intValue = int64(v)
	case int64:
		intValue = v
	case float64:
		if v != float64(int64(v)) {
			return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是整数", attr.Name))
		}
		intValue = int64(v)
	default:
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是整数类型", attr.Name))
	}

	// 验证最小值
	if attr.Validation.Min != nil && float64(intValue) < *attr.Validation.Min {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]不能小于%v", attr.Name, *attr.Validation.Min))
	}

	// 验证最大值
	if attr.Validation.Max != nil && float64(intValue) > *attr.Validation.Max {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]不能大于%v", attr.Name, *attr.Validation.Max))
	}

	return nil
}

// validateDateValue 验证日期值（秒级时间戳）
func (s *TemplateService) validateDateValue(attr model.Attribute, value interface{}) *herrors.HError {
	var timestamp int64
	switch v := value.(type) {
	case int:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case float64:
		timestamp = int64(v)
	default:
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是时间戳", attr.Name))
	}

	// 验证时间戳是否在合理范围内（1970-01-01 到 2100-12-31）
	minTimestamp := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	maxTimestamp := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC).Unix()

	if timestamp < minTimestamp || timestamp > maxTimestamp {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]时间戳超出有效范围", attr.Name))
	}

	// 验证是否为当天的开始时间（00:00:00）
	t := time.Unix(timestamp, 0).UTC()
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是当天的开始时间（00:00:00）", attr.Name))
	}

	return nil
}

// validateDateTimeValue 验证日期时间值（秒级时间戳）
func (s *TemplateService) validateDateTimeValue(attr model.Attribute, value interface{}) *herrors.HError {
	var timestamp int64
	switch v := value.(type) {
	case int:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case float64:
		timestamp = int64(v)
	default:
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是时间戳", attr.Name))
	}

	// 验证时间戳是否在合理范围内（1970-01-01 到 2100-12-31）
	minTimestamp := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	maxTimestamp := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC).Unix()

	if timestamp < minTimestamp || timestamp > maxTimestamp {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]时间戳超出有效范围", attr.Name))
	}

	return nil
}

// validateTimeValue 验证时间值
func (s *TemplateService) validateTimeValue(attr model.Attribute, value interface{}) *herrors.HError {
	strValue, ok := value.(string)
	if !ok {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是时间字符串", attr.Name))
	}

	// 尝试解析时间
	_, err := time.Parse("15:04:05", strValue)
	if err != nil {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]时间格式不正确，应为HH:mm:ss", attr.Name))
	}

	return nil
}

// validateSelectValue 验证选择值
func (s *TemplateService) validateSelectValue(attr model.Attribute, value interface{}) *herrors.HError {
	// 检查值是否在选项中
	for _, opt := range attr.Options {
		if opt.Value == value {
			return nil
		}
	}
	return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]的值不在可选范围内", attr.Name))
}

// validateBooleanValue 验证布尔值
func (s *TemplateService) validateBooleanValue(attr model.Attribute, value interface{}) *herrors.HError {
	_, ok := value.(bool)
	if !ok {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是布尔类型", attr.Name))
	}
	return nil
}

// validateWalletValue 验证钱包值
func (s *TemplateService) validateWalletValue(attr model.Attribute, value interface{}) *herrors.HError {
	// 钱包类型需要特殊处理，这里先实现基本的字符串验证
	strValue, ok := value.(string)
	if !ok {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]必须是字符串类型", attr.Name))
	}

	// 验证钱包地址格式（这里需要根据具体的钱包类型来实现）
	if len(strValue) < 10 {
		return template_err.TemplateValidationFailed(fmt.Errorf("字段[%s]钱包地址格式不正确", attr.Name))
	}

	return nil
}

// convertToModelAttribute 将值对象属性转换为领域模型属性
func convertToModelAttribute(attr valueobject.TemplateAttribute) model.Attribute {
	options := make([]model.Option, len(attr.Options))
	for i, opt := range attr.Options {
		options[i] = model.Option{
			Label: opt.Label,
			Value: opt.Value,
			Sort:  opt.Sort,
		}
	}

	return model.Attribute{
		Key:      attr.Key,
		Name:     attr.Name,
		Type:     attr.Type,
		Required: attr.Required,
		I18nKey:  attr.I18nKey,
		Options:  options,
		IsQuery:  attr.IsQuery,
		Default:  attr.Default,
		Validation: model.Validation{
			Min:     attr.Validation.Min,
			Max:     attr.Validation.Max,
			Pattern: attr.Validation.Pattern,
			Length:  attr.Validation.Length,
		},
		Description: attr.Description,
	}
}
