package excel

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// User 用户模型示例
type User struct {
	ID        int64   `json:"id" excel:"title:用户ID;width:10;sort:1"`
	Name      string  `json:"name" excel:"title:用户名;width:15;sort:2"`
	Age       int     `json:"age" excel:"title:年龄;width:8;sort:3"`
	Money     int64   `json:"money" excel:"title:余额;width:12;format:money;sort:4"`
	CreatedAt int64   `json:"created_at" excel:"title:创建时间;width:20;format:time;sort:5"`
	Score     float64 `json:"score" excel:"title:得分;width:10;format:percent;sort:6"`
}

// ExampleExport 导出示例
func ExampleExport() error {
	// 创建导出器
	exporter := NewExcelExporter("用户列表")

	// 方式1：从模型解析列
	exporter.ParseModelColumns(User{})

	// 方式2：手动设置列
	// exporter.SetColumns([]Column{
	//     {Title: "用户ID", Field: "id", Width: 10, Sort: 1},
	//     {Title: "用户名", Field: "name", Width: 15, Sort: 2},
	//     {Title: "年龄", Field: "age", Width: 8, Sort: 3},
	//     {Title: "余额", Field: "money", Width: 12, Format: "money", Sort: 4},
	//     {Title: "创建时间", Field: "created_at", Width: 20, Format: "time", Sort: 5},
	//     {Title: "得分", Field: "score", Width: 10, Format: "percent", Sort: 6},
	// })

	// 写入表头
	if err := exporter.WriteHeader(); err != nil {
		return err
	}

	// 准备数据
	users := []User{
		{
			ID:        1,
			Name:      "张三",
			Age:       25,
			Money:     10000, // 100.00元
			CreatedAt: utils.GetDateUnix(),
			Score:     0.85, // 85%
		},
		{
			ID:        2,
			Name:      "李四",
			Age:       30,
			Money:     20000, // 200.00元
			CreatedAt: utils.GetDateUnix(),
			Score:     0.92, // 92%
		},
	}

	// 写入数据
	if err := exporter.WriteRows(users); err != nil {
		return err
	}

	// 保存文件
	return exporter.Save("users.xlsx")
}

// ExampleExportWithCustomColumns 使用自定义列的导出示例
func ExampleExportWithCustomColumns() error {
	exporter := NewExcelExporter("自定义导出")

	// 设置自定义列
	exporter.SetColumns([]Column{
		{Title: "序号", Field: "index", Width: 8, Sort: 1},
		{Title: "名称", Field: "name", Width: 15, Sort: 2},
		{Title: "数值", Field: "value", Width: 12, Format: "money", Sort: 3},
	})

	// 写入表头
	if err := exporter.WriteHeader(); err != nil {
		return err
	}

	// 准备数据
	type CustomData struct {
		Index int    `json:"index"`
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	data := []CustomData{
		{Index: 1, Name: "项目A", Value: 1000},
		{Index: 2, Name: "项目B", Value: 2000},
		{Index: 3, Name: "项目C", Value: 3000},
	}

	// 写入数据
	if err := exporter.WriteRows(data); err != nil {
		return err
	}

	// 保存文件
	return exporter.Save("custom.xlsx")
}
