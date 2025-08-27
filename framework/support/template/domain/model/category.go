package model

import (
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// Category 模板分类领域模型
type Category struct {
	ID          string
	Name        string
	Code        string
	Description string
	Sort        int
	Status      int
	CreatedAt   int64
	UpdatedAt   int64
}

// NewCategory 创建模板分类
func NewCategory(name, code, description string) *Category {
	return &Category{
		Name:        name,
		Code:        code,
		Description: description,
		Sort:        0,
		Status:      1,
		CreatedAt:   utils.GetDateUnix(),
		UpdatedAt:   utils.GetDateUnix(),
	}
}

// ToJSON 转换为JSON字符串
func (c *Category) ToJSON() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析
func (c *Category) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), c)
}
