package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ExportType 导出类型
type ExportType int

const (
	ExportTypeCurrentPage ExportType = iota + 1 // 导出当前页
	ExportTypeQueryResult                       // 导出条件查询结果
	ExportTypeAll                               // 导出所有数据
)

// Column 定义导出列
type Column struct {
	Title  string  `json:"title"`  // 列标题
	Field  string  `json:"field"`  // 字段名
	Width  float64 `json:"width"`  // 列宽
	Format string  `json:"format"` // 格式化函数名
	Sort   int     `json:"sort"`   // 排序
}

// ExcelExporter Excel导出器
type ExcelExporter struct {
	file     *excelize.File
	sheet    string
	columns  []Column
	rowIndex int
}

// NewExcelExporter 创建Excel导出器
func NewExcelExporter(sheet string) *ExcelExporter {
	file := excelize.NewFile()
	// 创建新的 sheet
	file.NewSheet(sheet)
	// 删除默认的 Sheet1
	file.DeleteSheet("Sheet1")

	return &ExcelExporter{
		file:     file,
		sheet:    sheet,
		rowIndex: 1,
	}
}

// SetColumns 设置列
func (e *ExcelExporter) SetColumns(columns []Column) {
	e.columns = columns
}

// ParseModelColumns 从模型解析列
func (e *ExcelExporter) ParseModelColumns(model interface{}) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var columns []Column
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		excelTag := field.Tag.Get("excel")
		if excelTag == "" {
			continue
		}

		// 解析excel标签
		tags := strings.Split(excelTag, ";")
		column := Column{
			Field: field.Name,
			Title: field.Name,
			Width: 15,
		}

		// 优先使用json标签作为字段名
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			// 处理json标签中的omitempty等选项
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName != "" && jsonName != "-" {
				column.Field = jsonName
			}
		}

		for _, tag := range tags {
			kv := strings.Split(tag, ":")
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "title":
				column.Title = kv[1]
			case "width":
				s := kv[1]
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					column.Width = f
				}
			case "format":
				column.Format = kv[1]
			case "sort":
				f, err := strconv.Atoi(kv[1])
				if err == nil {
					column.Sort = f
				}
			}
		}
		columns = append(columns, column)
	}

	// 按sort排序
	for i := 0; i < len(columns); i++ {
		for j := i + 1; j < len(columns); j++ {
			if columns[i].Sort > columns[j].Sort {
				columns[i], columns[j] = columns[j], columns[i]
			}
		}
	}

	e.columns = columns
}

// WriteHeader 写入表头
func (e *ExcelExporter) WriteHeader() error {
	// 设置列宽
	for i, col := range e.columns {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}
		if err := e.file.SetColWidth(e.sheet, colName, colName, col.Width); err != nil {
			return err
		}
		// 写入标题
		cell := fmt.Sprintf("%s%d", colName, e.rowIndex)
		if err := e.file.SetCellValue(e.sheet, cell, col.Title); err != nil {
			return err
		}
	}
	e.rowIndex++
	return nil
}

// WriteRow 写入一行数据
func (e *ExcelExporter) WriteRow(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i, col := range e.columns {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}
		cell := fmt.Sprintf("%s%d", colName, e.rowIndex)

		// 获取字段值
		field := v.FieldByName(col.Field)
		if !field.IsValid() {
			// 如果通过字段名没找到，尝试通过json标签查找
			t := v.Type()
			for j := 0; j < t.NumField(); j++ {
				structField := t.Field(j)
				if jsonTag := structField.Tag.Get("json"); jsonTag != "" {
					jsonName := strings.Split(jsonTag, ",")[0]
					if jsonName == col.Field {
						field = v.Field(j)
						break
					}
				}
			}
		}

		if !field.IsValid() {
			continue
		}

		// 格式化值
		value := e.formatValue(field, col.Format)
		if err := e.file.SetCellValue(e.sheet, cell, value); err != nil {
			return err
		}
	}
	e.rowIndex++
	return nil
}

// WriteRows 写入多行数据
func (e *ExcelExporter) WriteRows(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if err := e.WriteRow(v.Index(i).Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Save 保存文件
func (e *ExcelExporter) Save(filename string) error {
	return e.file.SaveAs(filename)
}

// SaveAsBytes 保存为字节数组
func (e *ExcelExporter) SaveAsBytes() ([]byte, error) {
	// 创建内存缓冲区
	buffer, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("写入缓冲区失败: %w", err)
	}
	return buffer.Bytes(), nil
}

// formatValue 格式化值
func (e *ExcelExporter) formatValue(v reflect.Value, format string) interface{} {
	if !v.IsValid() {
		return nil
	}

	// 根据格式化函数处理值
	switch format {
	case "time":
		if v.Kind() == reflect.Int64 {
			return time.Unix(v.Int(), 0).Format("2006-01-02 15:04:05")
		}
	case "date":
		if v.Kind() == reflect.Int64 {
			return time.Unix(v.Int(), 0).Format("2006-01-02")
		}
	case "money":
		if v.Kind() == reflect.Int64 {
			return fmt.Sprintf("%.2f", float64(v.Int())/100)
		}
	case "percent":
		if v.Kind() == reflect.Float64 {
			return fmt.Sprintf("%.2f%%", v.Float()*100)
		}
	}

	// 默认返回原始值
	return v.Interface()
}

// GetColumns 获取当前设置的列配置
func (e *ExcelExporter) GetColumns() []Column {
	return e.columns
}
