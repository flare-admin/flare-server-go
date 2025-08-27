package dto

// ConfigDTO 配置数据传输对象
type ConfigDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Group       string `json:"group"`
	Description string `json:"description"`
	I18nKey     string `json:"i18n_key"`
	IsSystem    bool   `json:"is_system"`
	IsEnabled   bool   `json:"is_enabled"`
	Sort        int    `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// ConfigGroupDTO 配置分组数据传输对象
type ConfigGroupDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	I18nKey     string `json:"i18n_key"`
	IsSystem    bool   `json:"is_system"`
	IsEnabled   bool   `json:"is_enabled"`
	Sort        int    `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// ConfigValueDTO 配置值数据传输对象
type ConfigValueDTO struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ConfigValueMapDTO 配置值映射数据传输对象
type ConfigValueMapDTO struct {
	Values map[string]interface{} `json:"values"`
}

// ConfigQueryDTO 配置查询数据传输对象
type ConfigQueryDTO struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Type      string `json:"type"`
	Group     string `json:"group"`
	IsSystem  *bool  `json:"is_system"`
	IsEnabled *bool  `json:"is_enabled"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// ConfigGroupQueryDTO 配置分组查询数据传输对象
type ConfigGroupQueryDTO struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	IsSystem  *bool  `json:"is_system"`
	IsEnabled *bool  `json:"is_enabled"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}
