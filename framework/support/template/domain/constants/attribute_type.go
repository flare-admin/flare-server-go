package constants

// AttributeType 属性类型常量
const (
	// AttributeTypeWallet 钱包类型
	AttributeTypeWallet = "wallet"
	// AttributeTypeString 字符串类型
	AttributeTypeString = "string"
	// AttributeTypeNumber 数字类型
	AttributeTypeNumber = "number"
	// AttributeTypeInteger 整数类型
	AttributeTypeInteger = "integer"
	// AttributeTypeSelect 选择类型
	AttributeTypeSelect = "select"
	// AttributeTypeMoney 金额类型
	AttributeTypeMoney = "money"
	// AttributeTypeBoolean 布尔类型
	AttributeTypeBoolean = "boolean"
	// AttributeTypeDate 日期类型
	AttributeTypeDate = "date"
	// AttributeTypeDateTime 日期时间类型
	AttributeTypeDateTime = "datetime"
	// AttributeTypeTime 时间类型
	AttributeTypeTime = "time"
	// AttributeTypeSwitch 开关类型
	AttributeTypeSwitch = "switch"
	// AttributeTypeTextarea 文本域类型
	AttributeTypeTextarea = "textarea"
)

// AttributeTypeOptions 属性类型选项
var AttributeTypeOptions = []struct {
	Label string
	Value string
}{
	{Label: "attribute.type.wallet", Value: AttributeTypeWallet},
	{Label: "attribute.type.string", Value: AttributeTypeString},
	{Label: "attribute.type.number", Value: AttributeTypeNumber},
	{Label: "attribute.type.integer", Value: AttributeTypeInteger},
	{Label: "attribute.type.select", Value: AttributeTypeSelect},
	{Label: "attribute.type.money", Value: AttributeTypeMoney},
	{Label: "attribute.type.boolean", Value: AttributeTypeBoolean},
	{Label: "attribute.type.date", Value: AttributeTypeDate},
	{Label: "attribute.type.datetime", Value: AttributeTypeDateTime},
	{Label: "attribute.type.time", Value: AttributeTypeTime},
	{Label: "attribute.type.switch", Value: AttributeTypeSwitch},
	{Label: "attribute.type.textarea", Value: AttributeTypeTextarea},
}

// IsValidAttributeType 检查属性类型是否有效
func IsValidAttributeType(attrType string) bool {
	for _, option := range AttributeTypeOptions {
		if option.Value == attrType {
			return true
		}
	}
	return false
}
