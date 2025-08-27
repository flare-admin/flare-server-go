package dictionary

import (
	"errors"
	"sync"
)

// Category 表示字典分类
type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Option 表示字典选项
type Option struct {
	ID         string `json:"id"`
	CategoryID string `json:"category_id"`
	Code       string `json:"code"`   // 选项编码
	Value      string `json:"value"`  // 选项值
	Sort       int    `json:"sort"`   // 排序号
	Status     int    `json:"status"` // 状态：1-启用，0-禁用
	Remark     string `json:"remark"` // 备注
}

// Dictionary 表示字典管理器
type Dictionary struct {
	categories map[string]Category
	options    map[string][]Option
	mu         sync.RWMutex
}

// New 创建字典实例
func New() *Dictionary {
	return &Dictionary{
		categories: make(map[string]Category),
		options:    make(map[string][]Option),
	}
}

// AddCategory 添加分类
func (d *Dictionary) AddCategory(category Category) error {
	if category.ID == "" || category.Name == "" {
		return errors.New("分类ID和名称不能为空")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.categories[category.ID]; exists {
		return errors.New("分类ID已存在")
	}

	d.categories[category.ID] = category
	return nil
}

// GetCategory 获取分类
func (d *Dictionary) GetCategory(id string) (Category, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	category, exists := d.categories[id]
	if !exists {
		return Category{}, errors.New("分类不存在")
	}
	return category, nil
}

// ListCategories 列出所有分类
func (d *Dictionary) ListCategories() []Category {
	d.mu.RLock()
	defer d.mu.RUnlock()

	categories := make([]Category, 0, len(d.categories))
	for _, category := range d.categories {
		categories = append(categories, category)
	}
	return categories
}

// AddOption 添加选项
func (d *Dictionary) AddOption(option Option) error {
	if option.ID == "" || option.CategoryID == "" {
		return errors.New("选项ID和分类ID不能为空")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.categories[option.CategoryID]; !exists {
		return errors.New("分类不存在")
	}

	// 检查选项ID是否重复
	options := d.options[option.CategoryID]
	for _, opt := range options {
		if opt.ID == option.ID {
			return errors.New("选项ID已存在")
		}
	}

	d.options[option.CategoryID] = append(d.options[option.CategoryID], option)
	return nil
}

// GetOptions 获取分类下的所有选项
func (d *Dictionary) GetOptions(categoryID string) ([]Option, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if _, exists := d.categories[categoryID]; !exists {
		return nil, errors.New("分类不存在")
	}

	options := d.options[categoryID]
	result := make([]Option, len(options))
	copy(result, options)
	return result, nil
}

// UpdateOption 更新选项
func (d *Dictionary) UpdateOption(option Option) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	options := d.options[option.CategoryID]
	for i, opt := range options {
		if opt.ID == option.ID {
			options[i] = option
			return nil
		}
	}
	return errors.New("选项不存在")
}

// DeleteOption 删除选项
func (d *Dictionary) DeleteOption(categoryID, optionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	options := d.options[categoryID]
	for i, opt := range options {
		if opt.ID == optionID {
			// 从切片中删除元素
			d.options[categoryID] = append(options[:i], options[i+1:]...)
			return nil
		}
	}
	return errors.New("选项不存在")
}
