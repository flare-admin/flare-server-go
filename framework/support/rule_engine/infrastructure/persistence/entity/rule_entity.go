package entity

import "github.com/flare-admin/flare-server-go/framework/pkg/database"

// Rule 规则实体
type Rule struct {
	database.BaseModel
	ID              string `gorm:"primarykey"`
	Code            string `gorm:"size:100;not null;uniqueIndex;comment:规则编码"`
	Name            string `gorm:"size:100;not null;comment:规则名称"`
	Description     string `gorm:"size:500;comment:规则描述"`
	CategoryID      string `gorm:"not null;index;comment:分类ID"`
	TemplateID      string `gorm:"index;comment:模板ID（可选）"`
	Type            string `gorm:"size:50;not null;comment:规则类型：condition(条件规则) lua(lua脚本规则) formula(公式规则)"`
	Version         string `gorm:"size:20;not null;default:'1.0.0';comment:规则版本"`
	Status          int    `gorm:"not null;default:1;comment:状态：1-启用 2-禁用"`
	Scope           string `gorm:"size:50;not null;default:'global';comment:作用域：global(全局) product(商品) user(用户) order(订单) withdraw(提现) declare(申报) payment(支付)"`
	ScopeID         string `gorm:"size:100;comment:作用域ID（商品ID、用户ID、订单ID等）"`
	Triggers        string `gorm:"size:200;comment:触发动作列表(逗号分隔，如：create,update,delete,approve,reject,placeOrder,pay,withdraw,declare)"`
	ExecutionTiming string `gorm:"size:10;not null;default:'before';comment:执行时机：before(前置) after(后置) both(前后都执行)"`
	Conditions      string `gorm:"type:json;comment:条件表达式(JSON格式)"`
	LuaScript       string `gorm:"type:text;comment:Lua脚本代码"`
	Formula         string `gorm:"type:text;comment:计算公式"`
	FormulaVars     string `gorm:"type:json;comment:公式变量映射(JSON格式)"`
	Action          string `gorm:"size:50;not null;default:'allow';comment:触发动作：allow(允许) deny(拒绝) modify(修改) notify(通知) redirect(重定向)"`
	Priority        int32  `gorm:"not null;default:0;comment:优先级"`
	Sorting         int32  `gorm:"not null;default:0;comment:排序权重"`
	ExecuteCount    int64  `gorm:"not null;default:0;comment:执行次数"`
	SuccessCount    int64  `gorm:"not null;default:0;comment:成功次数"`
	LastExecuteAt   int64  `gorm:"not null;default:0;comment:最后执行时间"`
	TenantID        string `gorm:"size:50;comment:租户ID"`
}

func (Rule) TableName() string {
	return "rules"
}

// GetPrimaryKey 获取主键
func (Rule) GetPrimaryKey() string {
	return "id"
}
