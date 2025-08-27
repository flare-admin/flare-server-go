package lua_engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// DBOperationService 数据库操作服务
type DBOperationService struct {
	data database.IDataBase
}

// NewDBOperationService 创建数据库操作服务
func NewDBOperationService(data database.IDataBase) *DBOperationService {
	return &DBOperationService{
		data: data,
	}
}

// Insert 插入数据
func (s *DBOperationService) Insert(ctx context.Context, table string, data map[string]interface{}) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("插入数据不能为空")
	}

	db := s.data.DB(ctx)
	result := db.Table(table).Create(data)
	if result.Error != nil {
		return 0, fmt.Errorf("插入数据失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// Update 更新数据
// table: 表名
// data: 要更新的数据
// whereSQL: WHERE条件SQL语句，支持?占位符
// whereArgs: WHERE条件参数，按顺序对应SQL中的?占位符
func (s *DBOperationService) Update(ctx context.Context, table string, data map[string]interface{}, whereSQL string, whereArgs ...interface{}) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("更新数据不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 构建WHERE条件
	if whereSQL != "" {
		query = query.Where(whereSQL, whereArgs...)
	}

	result := query.Updates(data)
	if result.Error != nil {
		return 0, fmt.Errorf("更新数据失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// UpdateWithMap 使用map条件更新数据（向后兼容）
func (s *DBOperationService) UpdateWithMap(ctx context.Context, table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("更新数据不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 构建WHERE条件
	if len(where) > 0 {
		for field, value := range where {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	result := query.Updates(data)
	if result.Error != nil {
		return 0, fmt.Errorf("更新数据失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// Delete 删除数据
// table: 表名
// whereSQL: WHERE条件SQL语句，支持?占位符
// whereArgs: WHERE条件参数，按顺序对应SQL中的?占位符
func (s *DBOperationService) Delete(ctx context.Context, table string, whereSQL string, whereArgs ...interface{}) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 构建WHERE条件
	if whereSQL != "" {
		query = query.Where(whereSQL, whereArgs...)
	}

	result := query.Delete(&map[string]interface{}{})
	if result.Error != nil {
		return 0, fmt.Errorf("删除数据失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// DeleteWithMap 使用map条件删除数据（向后兼容）
func (s *DBOperationService) DeleteWithMap(ctx context.Context, table string, where map[string]interface{}) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 构建WHERE条件
	if len(where) > 0 {
		for field, value := range where {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	result := query.Delete(&map[string]interface{}{})
	if result.Error != nil {
		return 0, fmt.Errorf("删除数据失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// Query 查询数据
func (s *DBOperationService) Query(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	if sql == "" {
		return nil, fmt.Errorf("SQL语句不能为空")
	}

	db := s.data.DB(ctx)
	var results []map[string]interface{}

	result := db.Raw(sql, args...).Scan(&results)
	if result.Error != nil {
		return nil, fmt.Errorf("查询数据失败: %v", result.Error)
	}

	return results, nil
}

// QueryOne 查询单条数据
func (s *DBOperationService) QueryOne(ctx context.Context, sql string, args ...interface{}) (map[string]interface{}, error) {
	if sql == "" {
		return nil, fmt.Errorf("SQL语句不能为空")
	}

	db := s.data.DB(ctx)
	var result map[string]interface{}

	err := db.Raw(sql, args...).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询单条数据失败: %v", err)
	}

	return result, nil
}

// Execute 执行SQL语句
func (s *DBOperationService) Execute(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, fmt.Errorf("SQL语句不能为空")
	}

	db := s.data.DB(ctx)
	result := db.Exec(sql, args...)
	if result.Error != nil {
		return 0, fmt.Errorf("执行SQL失败: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// QueryWithBuilder 使用QueryBuilder查询数据
func (s *DBOperationService) QueryWithBuilder(ctx context.Context, table string, builder *db_query.QueryBuilder) ([]map[string]interface{}, error) {
	if table == "" {
		return nil, fmt.Errorf("表名不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 应用QueryBuilder的条件
	if builder != nil {
		query = builder.Build(query)
	}

	var results []map[string]interface{}
	err := query.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}

	return results, nil
}

// CountWithBuilder 使用QueryBuilder统计数量
func (s *DBOperationService) CountWithBuilder(ctx context.Context, table string, builder *db_query.QueryBuilder) (int64, error) {
	if table == "" {
		return 0, fmt.Errorf("表名不能为空")
	}

	db := s.data.DB(ctx)
	query := db.Table(table)

	// 应用QueryBuilder的条件（不包含分页和排序）
	if builder != nil {
		for _, cond := range builder.GetConditions() {
			if cond.IsRaw {
				query = query.Where(cond.RawSQL, cond.RawArgs...)
			} else {
				switch cond.Operator {
				case db_query.IsNull:
					query = query.Where(fmt.Sprintf("%s IS NULL", cond.Field))
				case db_query.IsNotNull:
					query = query.Where(fmt.Sprintf("%s IS NOT NULL", cond.Field))
				case db_query.In:
					query = query.Where(fmt.Sprintf("%s IN ?", cond.Field), cond.Value)
				case db_query.NotIn:
					query = query.Where(fmt.Sprintf("%s NOT IN ?", cond.Field), cond.Value)
				case db_query.Like:
					query = query.Where(fmt.Sprintf("%s LIKE ?", cond.Field), cond.Value)
				default:
					query = query.Where(fmt.Sprintf("%s %s ?", cond.Field, cond.Operator), cond.Value)
				}
			}
		}
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计数量失败: %v", err)
	}

	return count, nil
}

// BuildSQL 构建SQL语句
func (s *DBOperationService) BuildSQL(table string, builder *db_query.QueryBuilder) (string, []interface{}, error) {
	if table == "" {
		return "", nil, fmt.Errorf("表名不能为空")
	}

	var sql strings.Builder
	var args []interface{}

	// SELECT 子句
	sql.WriteString("SELECT * FROM " + table)

	// WHERE 子句
	if builder != nil {
		where, whereArgs := builder.BuildWhere()
		if where != "" {
			sql.WriteString(" WHERE " + where)
			args = append(args, whereArgs...)
		}

		// ORDER BY 子句
		orderBy := builder.BuildOrderBy()
		if orderBy != "" {
			sql.WriteString(" ORDER BY " + orderBy)
		}

		// LIMIT 子句
		limit, limitArgs := builder.BuildLimit()
		if limit != "" {
			sql.WriteString(" " + limit)
			args = append(args, limitArgs)
		}
	}

	return sql.String(), args, nil
}

// Transaction 事务操作
func (s *DBOperationService) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.data.InTx(ctx, fn)
}

// IndependentTransaction 独立事务操作
func (s *DBOperationService) IndependentTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.data.InIndependentTx(ctx, fn)
}
