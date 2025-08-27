# WhereDynamicQuery 复杂场景使用示例

## 1. 基础字段与 JSON 字段混合查询

```go
// 示例1：基础字段与 JSON 字段混合查询
quer := db_query.NewQueryBuilder()

// 基础字段查询
quer.Where("status", db_query.Eq, 1)
quer.Where("created_at", db_query.Gte, time.Now().AddDate(0, -1, 0))

// JSON 字段动态查询
quer.WhereDynamicQuery(map[string]interface{}{
    "name": "张三",
    "age": 25,
    "address.city": "北京",
    "address.district": "海淀区",
}, "content::jsonb -> '", "'")
```

## 2. 普通列查询

```go
// 示例2：普通列查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "user_name": "张三",
    "age": 25,
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "status": 1,
}, "", "")  // 不添加前缀和后缀，直接查询普通列
```

## 3. 混合列查询

```go
// 示例3：混合列查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "user_name": "张三",
    "age": 25,
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "address.city": "北京",
    "address.district": "海淀区",
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 4. 多表关联查询

```go
// 示例4：多表关联查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "d.name": "技术部",
    "d.code": "TECH",
    "u.name": "李四",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "department.manager.name": "李四",
    "department.employees[0].position": "工程师",
    "department.employees[0].skills": []string{"Go", "Python"},
}, "d.department_info::jsonb -> '", "'")  // JSON 字段查询
```

## 5. 条件组合查询

```go
// 示例5：条件组合查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "status": "active",
    "type": "project",
    "priority": "high",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "assignee.name": "王五",
    "tags[0]": "重要",
    "settings.notification": true,
    "settings.theme": "dark",
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 6. 时间范围查询

```go
// 示例6：时间范围查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "created_at": "2024-03-20",
    "updated_at": "2024-03-21",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "schedule.startTime": "09:00",
    "schedule.endTime": "18:00",
    "schedule.days": []string{"Monday", "Tuesday"},
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 7. 多语言内容查询

```go
// 示例7：多语言内容查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "language": "zh-CN",
    "country": "CN",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "title.zh": "标题",
    "title.en": "Title",
    "description.zh": "描述",
    "description.en": "Description",
    "keywords": []string{"标签1", "标签2"},
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 8. 权限相关查询

```go
// 示例8：权限相关查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "role": "admin",
    "permission_level": "high",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "permissions.roles[0]": "admin",
    "permissions.scopes[0]": "read",
    "permissions.resources[0].type": "document",
    "permissions.resources[0].action": "edit",
    "permissions.resources[0].conditions": map[string]interface{}{
        "department": "技术部",
        "level": "高级",
    },
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 9. 配置项查询

```go
// 示例9：配置项查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "config_type": "system",
    "config_status": "active",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "settings.theme": "dark",
    "settings.language": "zh-CN",
    "settings.notifications.email": true,
    "settings.notifications.sms": false,
    "settings.features": []string{"feature1", "feature2"},
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 10. 复杂的数据结构查询

```go
// 示例10：复杂的数据结构查询
quer := db_query.NewQueryBuilder()

// 查询普通列
quer.WhereDynamicQuery(map[string]interface{}{
    "project_id": "P001",
    "project_status": "active",
    "project_type": "development",
}, "", "")  // 普通列查询

// 查询 JSON 字段
quer.WhereDynamicQuery(map[string]interface{}{
    "project.name": "项目A",
    "project.members[0].role": "负责人",
    "project.tasks[0].status": "进行中",
    "project.settings.notification": true,
    "project.tags": []string{"重要", "紧急"},
    "project.metadata": map[string]interface{}{
        "created_by": "张三",
        "department": "技术部",
        "priority": "high",
    },
}, "content::jsonb -> '", "'")  // JSON 字段查询
```

## 使用注意事项

1. 性能优化
   - 合理使用索引
   - 避免过于复杂的嵌套查询
   - 考虑使用分页
   - 适当使用缓存

2. 数据安全
   - 注意 SQL 注入风险
   - 验证输入数据
   - 限制查询深度

3. 最佳实践
   - 保持查询结构清晰
   - 使用有意义的字段名
   - 适当添加注释
   - 考虑查询的可维护性

4. 错误处理
   - 处理空值情况
   - 处理类型转换错误
   - 处理查询超时
   - 记录错误日志

5. 列查询说明
   - 普通列查询：使用空字符串作为前缀和后缀
   - JSON 字段查询：使用 "content::jsonb -> '" 作为前缀， "'" 作为后缀
   - 可以同时使用普通列和 JSON 字段查询
   - 注意区分普通列和 JSON 字段的查询方式
