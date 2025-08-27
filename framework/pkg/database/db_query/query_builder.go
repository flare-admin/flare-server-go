package db_query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// 数据权限范围枚举
const (
	DataScopeAll      int8 = 1 // 全部数据
	DataScopeTenant   int8 = 2 // 本租户数据
	DataScopeDept     int8 = 3 // 本部门数据
	DataScopeDeptTree int8 = 4 // 本部门及以下数据
	DataScopeSelf     int8 = 5 // 仅本人数据
	DataScopeCustom   int8 = 6 // 自定义数据
)

// Operator 查询操作符
type Operator string

const (
	Eq        Operator = "="
	Neq       Operator = "!="
	Gt        Operator = ">"
	Gte       Operator = ">="
	Lt        Operator = "<"
	Lte       Operator = "<="
	Like      Operator = "LIKE"
	In        Operator = "IN"
	NotIn     Operator = "NOT IN"
	IsNull    Operator = "IS NULL"
	IsNotNull Operator = "IS NOT NULL"
)

// Condition 查询条件
type Condition struct {
	Field    string        // 字段名
	Operator Operator      // 操作符
	Value    interface{}   // 值
	RawSQL   string        // 原生SQL条件
	RawArgs  []interface{} // 原生SQL参数
	IsRaw    bool          // 是否为原生SQL条件
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	conditions []Condition
	orderBy    []string
	page       *Page
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		conditions: make([]Condition, 0),
		orderBy:    make([]string, 0),
	}
}

// Where 添加查询条件
func (qb *QueryBuilder) Where(field string, operator Operator, value interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
		IsRaw:    false,
	})
	return qb
}

// WhereRaw 添加原生SQL条件
func (qb *QueryBuilder) WhereRaw(sql string, args ...interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, Condition{
		RawSQL:  sql,
		RawArgs: args,
		IsRaw:   true,
	})
	return qb
}

// WhereIn 添加IN查询条件
func (qb *QueryBuilder) WhereIn(field string, values interface{}) *QueryBuilder {
	return qb.Where(field, In, values)
}

// WhereNotIn 添加NOT IN查询条件
func (qb *QueryBuilder) WhereNotIn(field string, values interface{}) *QueryBuilder {
	return qb.Where(field, NotIn, values)
}

// WhereLike 添加LIKE查询条件
func (qb *QueryBuilder) WhereLike(field string, value string) *QueryBuilder {
	return qb.Where(field, Like, "%"+value+"%")
}

// WhereEq 添加等于查询条件
func (qb *QueryBuilder) WhereEq(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Eq, value)
}

// WhereNeq 添加不等于查询条件
func (qb *QueryBuilder) WhereNeq(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Neq, value)
}

// WhereGt 添加大于查询条件
func (qb *QueryBuilder) WhereGt(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Gt, value)
}

// WhereGte 添加大于等于查询条件
func (qb *QueryBuilder) WhereGte(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Gte, value)
}

// WhereLt 添加小于查询条件
func (qb *QueryBuilder) WhereLt(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Lt, value)
}

// WhereLte 添加小于等于查询条件
func (qb *QueryBuilder) WhereLte(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, Lte, value)
}

// WhereIsNull 添加IS NULL查询条件
func (qb *QueryBuilder) WhereIsNull(field string) *QueryBuilder {
	return qb.Where(field, IsNull, nil)
}

// WhereIsNotNull 添加IS NOT NULL查询条件
func (qb *QueryBuilder) WhereIsNotNull(field string) *QueryBuilder {
	return qb.Where(field, IsNotNull, nil)
}

// OrderBy 添加排序
func (qb *QueryBuilder) OrderBy(field string, asc bool) *QueryBuilder {
	direction := "DESC"
	if asc {
		direction = "ASC"
	}
	qb.orderBy = append(qb.orderBy, fmt.Sprintf("%s %s", field, direction))
	return qb
}

// OrderByDESC 添加降序排序
func (qb *QueryBuilder) OrderByDESC(field string) *QueryBuilder {
	return qb.OrderBy(field, false)
}

// OrderByASC 添加升序排序
func (qb *QueryBuilder) OrderByASC(field string) *QueryBuilder {
	return qb.OrderBy(field, true)
}

// WithPage 设置分页
func (qb *QueryBuilder) WithPage(page *Page) *QueryBuilder {
	qb.page = page
	return qb
}

// BuildWhere 构建WHERE子句
func (qb *QueryBuilder) BuildWhere() (string, []interface{}) {
	if len(qb.conditions) == 0 {
		return "", nil
	}

	var (
		where  strings.Builder
		values []interface{}
	)

	for i, cond := range qb.conditions {
		if i > 0 {
			where.WriteString(" AND ")
		}

		if cond.IsRaw {
			// 处理原生SQL条件
			where.WriteString(cond.RawSQL)
			values = append(values, cond.RawArgs...)
		} else {
			// 处理标准条件
			switch cond.Operator {
			case IsNull, IsNotNull:
				where.WriteString(fmt.Sprintf("%s %s", cond.Field, cond.Operator))
			case In, NotIn:
				where.WriteString(fmt.Sprintf("%s %s (?)", cond.Field, cond.Operator))
				values = append(values, cond.Value)
			default:
				where.WriteString(fmt.Sprintf("%s %s ?", cond.Field, cond.Operator))
				values = append(values, cond.Value)
			}
		}
	}

	return where.String(), values
}

// BuildOrderBy 构建ORDER BY子句
func (qb *QueryBuilder) BuildOrderBy() string {
	if len(qb.orderBy) == 0 {
		return ""
	}
	return strings.Join(qb.orderBy, ", ")
	//return "ORDER BY " + strings.Join(qb.orderBy, ", ")
}

// BuildLimit 构建LIMIT子句
func (qb *QueryBuilder) BuildLimit() (string, []int) {
	if qb.page == nil || qb.page.NoUse {
		return "", nil
	}
	qb.page.Fix()
	return "LIMIT ?, ?", []int{qb.page.Offset(), qb.page.Limit()}
}

// Build 将查询条件应用到GORM的DB对象上
func (qb *QueryBuilder) Build(db *gorm.DB) *gorm.DB {
	// 1. 应用WHERE条件
	for _, cond := range qb.conditions {
		if cond.IsRaw {
			// 处理原生SQL条件
			db = db.Where(cond.RawSQL, cond.RawArgs...)
		} else {
			// 处理标准条件
			switch cond.Operator {
			case IsNull:
				db = db.Where(fmt.Sprintf("%s IS NULL", cond.Field))
			case IsNotNull:
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", cond.Field))
			case In:
				db = db.Where(fmt.Sprintf("%s IN ?", cond.Field), cond.Value)
			case NotIn:
				db = db.Where(fmt.Sprintf("%s NOT IN ?", cond.Field), cond.Value)
			case Like:
				db = db.Where(fmt.Sprintf("%s LIKE ?", cond.Field), cond.Value)
			default:
				db = db.Where(fmt.Sprintf("%s %s ?", cond.Field, cond.Operator), cond.Value)
			}
		}
	}

	// 2. 应用ORDER BY
	if len(qb.orderBy) > 0 {
		for _, order := range qb.orderBy {
			db = db.Order(order)
		}
	}

	// 3. 应用分页
	if qb.page != nil && !qb.page.NoUse {
		qb.page.Fix()
		db = db.Offset(qb.page.Offset()).Limit(qb.page.Limit())
	}

	return db
}

// GetConditions 获取查询条件
func (qb *QueryBuilder) GetConditions() []Condition {
	return qb.conditions
}

// WhereDynamicQuery 添加动态查询条件
// dynamicQuery: 动态查询参数，key为字段名，value为查询值
// fieldPrefix: 字段前缀，如 "content::jsonb ->> '"
// fieldSuffix: 字段后缀，如 "'"
func (qb *QueryBuilder) WhereDynamicQuery(dynamicQuery map[string]interface{}, fieldPrefix, fieldSuffix string) *QueryBuilder {
	if dynamicQuery == nil || len(dynamicQuery) == 0 {
		return qb
	}

	for key, value := range dynamicQuery {
		// 构建完整的字段名
		field := fmt.Sprintf("%s%s%s", fieldPrefix, key, fieldSuffix)

		// 根据 value 的类型构建不同的查询条件
		switch v := value.(type) {
		case string:
			// 字符串类型直接比较
			qb.Where(field, Like, "%"+v+"%")
		case int, int32, int64, float32, float64:
			// 数字类型转换为字符串后比较
			qb.Where(field, Eq, fmt.Sprintf("%v", v))
		case bool:
			// 布尔类型转换为字符串后比较
			qb.Where(field, Eq, fmt.Sprintf("%v", v))
		default:
			// 其他类型转换为字符串后比较
			qb.Where(field, Eq, fmt.Sprintf("%v", v))
		}
	}
	return qb
}
