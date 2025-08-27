package dictionary

import (
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := New()

	// 测试添加分类
	category := Category{
		ID:          "gender",
		Name:        "性别",
		Description: "用户性别选项",
	}

	err := dict.AddCategory(category)
	if err != nil {
		t.Errorf("添加分类失败: %v", err)
	}

	// 测试添加选项
	option1 := Option{
		ID:         "male",
		CategoryID: "gender",
		Code:       "1",
		Value:      "男",
		Sort:       1,
		Status:     1,
	}

	option2 := Option{
		ID:         "female",
		CategoryID: "gender",
		Code:       "2",
		Value:      "女",
		Sort:       2,
		Status:     1,
	}

	err = dict.AddOption(option1)
	if err != nil {
		t.Errorf("添加选项失败: %v", err)
	}

	err = dict.AddOption(option2)
	if err != nil {
		t.Errorf("添加选项失败: %v", err)
	}

	// 测试获取选项
	options, err := dict.GetOptions("gender")
	if err != nil {
		t.Errorf("获取选项失败: %v", err)
	}
	if len(options) != 2 {
		t.Errorf("期望获得2个选项，实际获得%d个", len(options))
	}

	// 测试更新选项
	option1.Value = "男性"
	err = dict.UpdateOption(option1)
	if err != nil {
		t.Errorf("更新选项失败: %v", err)
	}

	// 测试删除选项
	err = dict.DeleteOption("gender", "male")
	if err != nil {
		t.Errorf("删除选项失败: %v", err)
	}

	options, _ = dict.GetOptions("gender")
	if len(options) != 1 {
		t.Errorf("期望删除后剩余1个选项，实际剩余%d个", len(options))
	}
}
