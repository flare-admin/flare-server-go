# Lua引擎结构体支持功能

## 概述

Lua引擎现在支持从上下文中直接获取结构体，包括 `UserDTO` 和通用结构体。这个功能允许Lua脚本修改上下文中的结构体，并在执行完成后获取修改后的结构体。

## 功能特性

### 1. 支持的结构体类型

- **UserDTO**: 专门的用户DTO获取方法
- **通用结构体**: 支持任意结构体的获取和转换
- **嵌套结构体**: 支持嵌套路径的结构体获取
- **复杂结构体**: 支持包含map、slice等复杂字段的结构体

### 2. 获取方法

#### 专门方法（推荐）

```go
// 获取UserDTO
if userDTO, ok := result.GetUserDTO("user"); ok {
    fmt.Printf("用户ID: %s\n", userDTO.Id)
    fmt.Printf("用户昵称: %s\n", userDTO.Nickname)
    fmt.Printf("用户等级: %d\n", userDTO.Level)
}

// 获取嵌套UserDTO
if nestedUserDTO, ok := result.GetNestedUserDTO("nested.user"); ok {
    fmt.Printf("嵌套用户等级: %d\n", nestedUserDTO.Level)
}
```

#### 通用方法

```go
// 获取任意结构体
var userStruct dto.UserDTO
if result.GetContextStruct("user", &userStruct) {
    fmt.Printf("用户等级: %d\n", userStruct.Level)
}

// 获取嵌套任意结构体
var nestedUserStruct dto.UserDTO
if result.GetNestedContextStruct("nested.user", &nestedUserStruct) {
    fmt.Printf("嵌套用户等级: %d\n", nestedUserStruct.Level)
}
```

## 使用示例

### 1. 基本使用

```go
// 创建用户DTO
user := &dto.UserDTO{
    Id:             "user_001",
    Nickname:       "张三",
    Phone:          "13800138000",
    Level:          1,
    Status:         1,
    InvitationCode: "ABC123",
    Email:          "zhangsan@example.com",
    Age:            25,
    Gender:         2,
}

// 创建上下文数据
context := map[string]interface{}{
    "user": user,
    "nested": map[string]interface{}{
        "user": user,
    },
}

// 创建规则执行器
executor := lua_engine.NewRuleExecutor()

// 创建执行选项
opts := lua_engine.NewExecuteOptions().WithContext(context)

// Lua脚本：修改用户信息
script := `
    -- 修改用户等级
    set_object_property("user", "level", 3)
    
    -- 修改用户昵称
    set_object_property("user", "nickname", "李四")
    
    -- 修改用户状态
    set_object_property("user", "status", 2)
    
    -- 修改嵌套用户信息
    set_object_property("nested.user", "level", 4)
    set_object_property("nested.user", "nickname", "王五")
    
    -- 设置规则执行结果
    valid = true
    action = "update_user"
    variables = {
        old_level = 1,
        new_level = 3
    }
`

// 执行规则
result, err := executor.Execute(script, opts)
if err != nil {
    log.Fatalf("规则执行失败: %v", err)
}

// 获取修改后的UserDTO
if userDTO, ok := result.GetUserDTO("user"); ok {
    fmt.Printf("获取到的UserDTO: %+v\n", userDTO)
    fmt.Printf("用户等级: %d\n", userDTO.Level)        // 输出: 3
    fmt.Printf("用户昵称: %s\n", userDTO.Nickname)      // 输出: 李四
    fmt.Printf("用户状态: %d\n", userDTO.Status)        // 输出: 2
}
```

### 2. 通用结构体使用

```go
// 定义自定义结构体
type CustomStruct struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Age      int     `json:"age"`
    Score    float64 `json:"score"`
    IsActive bool    `json:"isActive"`
}

// 创建自定义结构体数据
customData := &CustomStruct{
    ID:       "custom_001",
    Name:     "自定义用户",
    Age:      30,
    Score:    95.5,
    IsActive: true,
}

// 创建上下文数据
context := map[string]interface{}{
    "custom": customData,
}

// 创建执行选项
opts := lua_engine.NewExecuteOptions().WithContext(context)

// Lua脚本：修改自定义结构体
script := `
    -- 修改自定义结构体
    set_object_property("custom", "name", "修改后的自定义用户")
    set_object_property("custom", "age", 35)
    set_object_property("custom", "score", 98.5)
    set_object_property("custom", "isActive", false)
    
    -- 设置规则执行结果
    valid = true
    action = "update_custom_struct"
    variables = {
        modified = true
    }
`

// 执行规则
result, err := executor.Execute(script, opts)
if err != nil {
    log.Fatalf("规则执行失败: %v", err)
}

// 获取修改后的自定义结构体
var customStruct CustomStruct
if result.GetContextStruct("custom", &customStruct) {
    fmt.Printf("获取到的CustomStruct: %+v\n", customStruct)
    fmt.Printf("自定义用户名称: %s\n", customStruct.Name)      // 输出: 修改后的自定义用户
    fmt.Printf("自定义用户年龄: %d\n", customStruct.Age)        // 输出: 35
    fmt.Printf("自定义用户分数: %.1f\n", customStruct.Score)    // 输出: 98.5
    fmt.Printf("自定义用户是否激活: %t\n", customStruct.IsActive) // 输出: false
}
```

### 3. 从map转换为结构体

```go
// 创建map数据
mapData := map[string]interface{}{
    "user": map[string]interface{}{
        "id":             "user_003",
        "nickname":       "赵六",
        "phone":          "13900139000",
        "level":          2,
        "status":         1,
        "invitationCode": "XYZ789",
        "email":          "zhaoliu@example.com",
        "age":            28,
        "gender":         1,
    },
}

// 创建执行选项
opts := lua_engine.NewExecuteOptions().WithContext(mapData)

// 简单的Lua脚本
script := `
    valid = true
    action = "test_map_conversion"
`

// 执行规则
result, err := executor.Execute(script, opts)
if err != nil {
    log.Fatalf("规则执行失败: %v", err)
}

// 从map转换为UserDTO
if userDTO, ok := result.GetUserDTO("user"); ok {
    fmt.Printf("从map转换的UserDTO: %+v\n", userDTO)
    fmt.Printf("用户ID: %s\n", userDTO.Id)           // 输出: user_003
    fmt.Printf("用户昵称: %s\n", userDTO.Nickname)     // 输出: 赵六
    fmt.Printf("用户等级: %d\n", userDTO.Level)        // 输出: 2
    fmt.Printf("用户年龄: %d\n", userDTO.Age)          // 输出: 28
    fmt.Printf("用户性别: %d\n", userDTO.Gender)       // 输出: 1
}
```

### 4. 复杂结构体使用

```go
// 定义复杂结构体
type ComplexStruct struct {
    ID       string                 `json:"id"`
    Name     string                 `json:"name"`
    Metadata map[string]interface{} `json:"metadata"`
    Tags     []string               `json:"tags"`
    Settings struct {
        Theme    string `json:"theme"`
        Language string `json:"language"`
    } `json:"settings"`
}

// 创建复杂结构体数据
complexData := map[string]interface{}{
    "id":   "complex_001",
    "name": "复杂结构体",
    "metadata": map[string]interface{}{
        "version": "1.0.0",
        "author":  "张三",
        "tags":    []string{"test", "demo"},
    },
    "tags": []string{"important", "urgent"},
    "settings": map[string]interface{}{
        "theme":    "dark",
        "language": "zh-CN",
    },
}

// 创建上下文数据
context := map[string]interface{}{
    "complex": complexData,
}

// 创建执行选项
opts := lua_engine.NewExecuteOptions().WithContext(context)

// Lua脚本：修改复杂结构体
script := `
    -- 修改复杂结构体
    set_object_property("complex", "name", "修改后的复杂结构体")
    set_nested_property("complex", "metadata.version", "2.0.0")
    set_nested_property("complex", "settings.theme", "light")
    
    -- 设置规则执行结果
    valid = true
    action = "update_complex_struct"
    variables = {
        modified = true
    }
`

// 执行规则
result, err := executor.Execute(script, opts)
if err != nil {
    log.Fatalf("规则执行失败: %v", err)
}

// 获取复杂结构体
var complexStruct ComplexStruct
if result.GetContextStruct("complex", &complexStruct) {
    fmt.Printf("获取到的ComplexStruct: %+v\n", complexStruct)
    fmt.Printf("复杂结构体名称: %s\n", complexStruct.Name)           // 输出: 修改后的复杂结构体
    fmt.Printf("复杂结构体主题: %s\n", complexStruct.Settings.Theme)   // 输出: light
    fmt.Printf("复杂结构体版本: %v\n", complexStruct.Metadata["version"]) // 输出: 2.0.0
    fmt.Printf("复杂结构体标签数量: %d\n", len(complexStruct.Tags))   // 输出: 2
}
```

## Lua脚本中的结构体修改

### 可用的Lua函数

1. **set_object_property(key, property, value)**: 设置对象属性
2. **set_nested_property(path, value)**: 设置嵌套属性
3. **set_context_value(key, value)**: 设置上下文值

### 示例

```lua
-- 修改用户等级
set_object_property("user", "level", 3)

-- 修改用户昵称
set_object_property("user", "nickname", "李四")

-- 修改用户状态
set_object_property("user", "status", 2)

-- 修改嵌套用户信息
set_object_property("nested.user", "level", 4)
set_object_property("nested.user", "nickname", "王五")

-- 修改嵌套属性
set_nested_property("complex.metadata.version", "2.0.0")
set_nested_property("complex.settings.theme", "light")

-- 设置规则执行结果
valid = true
action = "update_user"
variables = {
    old_level = 1,
    new_level = 3
}
```

## 注意事项

1. **结构体字段必须支持JSON序列化**: 所有需要转换的结构体字段都必须有 `json` 标签
2. **类型安全**: 使用专门的方法（如 `GetUserDTO`）比通用方法更安全
3. **错误处理**: 始终检查返回值，确保获取成功
4. **性能考虑**: 对于大量数据，建议使用专门的方法而不是通用方法
5. **嵌套路径**: 使用点号分隔的路径来访问嵌套结构体

## 支持的字段类型

- 基本类型: `string`, `int`, `int32`, `int64`, `float32`, `float64`, `bool`
- 复杂类型: `map[string]interface{}`, `[]string`, `[]int`, 等
- 嵌套结构体: 支持任意深度的嵌套结构体
- 指针类型: 支持结构体指针

## 错误处理

```go
// 检查获取是否成功
if userDTO, ok := result.GetUserDTO("user"); ok {
    // 成功获取，使用userDTO
    fmt.Printf("用户: %+v\n", userDTO)
} else {
    // 获取失败，处理错误
    fmt.Println("无法获取UserDTO")
}

// 检查通用方法获取是否成功
var userStruct dto.UserDTO
if result.GetContextStruct("user", &userStruct) {
    // 成功获取，使用userStruct
    fmt.Printf("用户: %+v\n", userStruct)
} else {
    // 获取失败，处理错误
    fmt.Println("无法获取UserDTO")
}
``` 