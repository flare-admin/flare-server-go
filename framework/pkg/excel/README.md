 # Excel 导出工具

这是一个通用的 Excel 导出工具，支持以下特性：

1. 支持从结构体标签自动解析导出列
2. 支持手动设置导出列
3. 支持自定义列宽
4. 支持自定义格式化函数
5. 支持列排序

## 安装依赖

```bash
go get github.com/xuri/excelize/v2
```

## 使用方法

### 1. 使用结构体标签导出

```go
// 定义结构体
type User struct {
    ID        int64     `excel:"title:用户ID;width:10;sort:1"`
    Name      string    `excel:"title:用户名;width:15;sort:2"`
    Age       int       `excel:"title:年龄;width:8;sort:3"`
    Money     int64     `excel:"title:余额;width:12;format:money;sort:4"`
    CreatedAt int64     `excel:"title:创建时间;width:20;format:time;sort:5"`
    Score     float64   `excel:"title:得分;width:10;format:percent;sort:6"`
}

// 导出数据
exporter := excel.NewExcelExporter("用户列表")
exporter.ParseModelColumns(User{})
exporter.WriteHeader()
exporter.WriteRows(users)
exporter.Save("users.xlsx")
```

### 2. 使用自定义列导出

```go
// 设置自定义列
exporter.SetColumns([]excel.Column{
    {Title: "序号", Field: "Index", Width: 8, Sort: 1},
    {Title: "名称", Field: "Name", Width: 15, Sort: 2},
    {Title: "数值", Field: "Value", Width: 12, Format: "money", Sort: 3},
})

// 导出数据
exporter.WriteHeader()
exporter.WriteRows(data)
exporter.Save("custom.xlsx")
```

## 标签说明

结构体标签格式：`excel:"key1:value1;key2:value2"`

支持的标签：

- `title`: 列标题
- `width`: 列宽
- `format`: 格式化函数
- `sort`: 排序序号

## 格式化函数

内置的格式化函数：

- `time`: 时间格式化 (2006-01-02 15:04:05)
- `date`: 日期格式化 (2006-01-02)
- `money`: 金额格式化 (除以100并保留2位小数)
- `percent`: 百分比格式化 (乘以100并添加%符号)

## 示例

查看 `example.go` 文件获取完整示例。