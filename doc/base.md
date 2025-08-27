## 主要业务
### 1、用户实名
- 奖励自己可以多个奖励
- 奖励上级奖励
- 按照比例奖励上级和上级等
### 2、用户注册
- 奖励自己可以多个奖励
- 奖励上级奖励
- 按照比例奖励上级和上级等
### 3、用户签到
- 奖励上级奖励
- 按照比例奖励上级和上级等
### 4、完成订单

## 任务中心

### 1. 功能概述

任务中心是一个灵活的任务管理系统，支持多种任务类型、奖励机制和返佣功能。系统可以处理普通任务和立即触发任务，并提供完整的任务生命周期管理。

### 2. 核心功能

#### 2.1 任务类型
- 普通任务：需要用户完成特定条件才能获得奖励
- 立即触发任务：如注册、签到、实名认证等，触发即完成

#### 2.2 任务分类
- 支持任务分类管理
- 每个分类可以设置不同的完成方式和重置规则
- 分类支持进度自动完成和可重复完成两种模式

#### 2.3 任务奖励
- 支持多种奖励类型（现金、积分、虚拟币等）
- 奖励发放方式：
  - 自动发放：任务完成时自动发放奖励
  - 手动领取：用户需要手动领取奖励
- 支持多级返佣机制：
  - 可配置不同等级的返佣金额
  - 返佣金额可来自任务奖励或自定义金额
  - 支持多种返佣来源类型

#### 2.4 任务条件
- 支持多种条件类型：
  - 数量达标
  - 金额达标
  - 时间达标
  - 特定事件
- 支持多种比较操作符（>=, <=, =, >, <）

#### 2.5 任务状态
- 未开始
- 进行中
- 已完成
- 已领取奖励
- 已重置
- 失败

### 3. 技术实现

#### 3.2 核心接口
```go
// 任务管理接口
type ITaskManageService interface {
    AddTask(ctx context.Context, req *dto.AddTask) herrors.Herr
    UpdateTask(ctx context.Context, req *dto.UpdateTask) herrors.Herr
    DeleteTask(ctx context.Context, id string) herrors.Herr
    GetTaskListByCategory(ctx context.Context, req *dto.GetTaskListReq) ([]*dto.TaskInfo, int64, herrors.Herr)
    GetTaskDetail(ctx context.Context, id string) (*dto.TaskInfo, herrors.Herr)
}

// 用户任务接口
type IUserTaskService interface {
    CompleteTask(ctx context.Context, req *dto.CompleteTask) herrors.Herr
    ClaimTaskReward(ctx context.Context, req *dto.ClaimTaskReward) herrors.Herr
    GetTaskProgress(ctx context.Context, req *dto.TaskAndUserIdQuery) ([]*dto.TaskConditionProgress, herrors.Herr)
    GetCategoryTaskStatus(ctx context.Context, req *dto.GetUserTaskStatusReq) ([]*dto.TaskStatusInfo, herrors.Herr)
}
```

### 4. 业务流程

#### 4.1 任务创建流程
1. 创建任务基本信息
2. 配置任务奖励
3. 配置任务返佣（可选）
4. 设置任务条件
5. 设置任务时间范围

#### 4.2 任务完成流程
1. 检查任务状态和有效性
2. 验证任务条件是否满足
3. 处理任务奖励发放
4. 处理任务返佣（如果配置了返佣）
5. 更新任务状态

#### 4.3 立即触发任务流程
1. 触发任务事件
2. 直接处理任务奖励
3. 处理任务返佣
4. 记录任务完成状态

### 5. 注意事项

#### 5.1 性能考虑
- 使用缓存优化任务状态查询
- 批量处理任务奖励和返佣
- 合理使用数据库索引

#### 5.2 安全考虑
- 任务奖励和返佣金额的精确计算
- 防止重复发放奖励
- 事务处理确保数据一致性

#### 5.3 扩展性考虑
- 支持自定义任务类型
- 支持自定义奖励类型
- 支持自定义返佣规则

### 6. 后续优化方向

1. 任务模板功能
2. 任务推荐系统
3. 任务数据分析
4. 任务活动管理
5. 任务通知系统