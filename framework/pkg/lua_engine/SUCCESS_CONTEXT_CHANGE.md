# Success函数Context变化机制

## 概述

本文档详细说明了success函数中variables设置后context的变化机制，以及如何验证这些变化。

## 执行流程

### 1. 初始Context注入
```go
// 在executor.go中
contextTable := L.NewTable()
e.injectContextData(L, contextTable, opts.Context)
L.SetGlobal("context", contextTable)
```

### 2. Success函数执行
```lua
-- 在Lua脚本中
success("action_name", {
    variable1 = "value1",
    variable2 = "value2",
    nested = {
        key = "value"
    }
})
```

### 3. Variables设置到Context
```go
// 在helpers.go的success函数中
contextTable := e.getOrCreateContextTable(L)
variables.ForEach(func(k, v lua.LValue) {
    contextTable.RawSetString(k.String(), v)
})
```

### 4. 提取修改后的Context
```go
// 在executor.go中
modifiedContext := e.extractModifiedContext(L, opts.Context)
result.Context = modifiedContext
```

## 验证机制

### 新增的辅助函数

#### 1. `get_context_size()`
```lua
-- 获取context表的大小
local size = get_context_size()
print("Context大小:", size)
```

#### 2. `print_context_keys()`
```lua
-- 打印context中的所有键值对
print_context_keys()
```

#### 3. `verify_context_value(key)`
```lua
-- 验证指定键的值是否存在
local value, exists = verify_context_value("key_name")
print("值:", value, "存在:", exists)
```

## 使用示例

### 基本使用
```lua
-- 设置variables
success("update_user", {
    old_age = 25,
    new_age = 30,
    reason = "age_update",
    processed_at = now(),
    processor = "rule_engine"
})

-- 验证变化
print("Context大小:", get_context_size())
print_context_keys()

local old_age, exists = verify_context_value("old_age")
print("old_age:", old_age, "exists:", exists)
```

### 复杂对象设置
```lua
-- 设置复杂对象
success("complex_update", {
    user = {
        id = "user_002",
        name = "李四",
        age = 30,
        level = 2
    },
    order = {
        status = "completed",
        processed_at = now()
    },
    metadata = {
        source = "rule_engine",
        version = "1.0.0"
    }
})
```

## 验证测试

### 运行演示
```go
// 运行基本演示
lua_engine.SuccessContextDemo()

// 运行详细测试
lua_engine.DemonstrateSuccessFunctionContextChange()
```

### 测试覆盖
- ✅ Variables正确设置到context中
- ✅ Variables与原始context数据合并
- ✅ 相同键的variables覆盖原始值
- ✅ 新的variables添加到context中
- ✅ 修改后的context在执行结果中返回
- ✅ 原始数据保持不变

## 技术细节

### Context变化流程
1. **初始注入**: 通过`injectContextData`注入初始context数据
2. **Lua执行**: Lua脚本执行，包括success函数调用
3. **Variables设置**: success函数将variables设置到context中
4. **Context提取**: 通过`extractModifiedContext`提取修改后的context
5. **结果返回**: 修改后的context作为执行结果返回

### 数据合并规则
- **新增变量**: 直接添加到context中
- **覆盖变量**: 相同键的variables会覆盖原始context中的值
- **嵌套对象**: 支持设置复杂的嵌套对象结构
- **原始数据**: 未被覆盖的原始数据保持不变

### 性能考虑
- **内存效率**: 使用Lua表直接操作，避免不必要的转换
- **执行速度**: 变量设置操作时间复杂度为O(1)
- **内存管理**: 自动处理Lua表的生命周期

## 最佳实践

### 1. 变量命名
```lua
-- 使用有意义的变量名
success("user_update", {
    old_user_level = 1,
    new_user_level = 2,
    upgrade_reason = "vip_promotion",
    processed_timestamp = now()
})
```

### 2. 结构化数据
```lua
-- 使用结构化数据便于后续处理
success("order_processing", {
    order_status = "completed",
    processing_result = {
        success = true,
        message = "订单处理成功",
        timestamp = now()
    },
    metadata = {
        processor = "rule_engine",
        version = "1.0.0"
    }
})
```

### 3. 验证变化
```lua
-- 在设置后验证变化
success("test_action", {
    test_var = "test_value"
})

-- 验证设置是否成功
local value, exists = verify_context_value("test_var")
if exists then
    print("变量设置成功:", value)
else
    print("变量设置失败")
end
```

## 故障排除

### 常见问题

#### 1. Variables未设置到Context
**原因**: success函数调用失败或context表不存在
**解决**: 检查Lua脚本语法和success函数调用

#### 2. Context变化未反映到结果
**原因**: extractModifiedContext方法未正确提取
**解决**: 检查context表的转换逻辑

#### 3. 原始数据被意外修改
**原因**: variables覆盖了不应该覆盖的键
**解决**: 使用更具体的变量名避免冲突

### 调试技巧
```lua
-- 在关键位置添加调试信息
print("Context大小:", get_context_size())
print_context_keys()

-- 验证特定变量
local value, exists = verify_context_value("target_key")
print("目标变量:", value, "存在:", exists)
```

## 总结

Success函数中的variables设置机制确保了：

1. **数据完整性**: Variables正确设置到context中
2. **数据合并**: 与原始context数据正确合并
3. **结果传递**: 修改后的context在执行结果中返回
4. **验证能力**: 提供多种方式验证context变化
5. **调试支持**: 丰富的调试和验证工具

这个机制为规则引擎提供了强大的数据传递和修改能力，使得复杂的业务逻辑能够通过简单的Lua脚本实现。 