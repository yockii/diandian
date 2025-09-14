package model

type Setting struct {
	Base
	Key         string  `json:"key,omitempty" gorm:"size:200;uniqueIndex"`
	Value       *string `json:"value,omitempty" gorm:"type:text"`
	GroupName   string  `json:"group_name,omitempty" gorm:"size:100;index"`
	Name        string  `json:"name,omitempty" gorm:"size:50"`
	Desc        string  `json:"desc,omitempty" gorm:"size:200"`
	OrderNum    int     `json:"order_num,omitempty" gorm:"default:0"`
	Showable    *bool   `json:"showable,omitempty" gorm:"default:true"` // 是否在设置界面展示
	SettingType string  `json:"setting_type,omitempty" gorm:"size:20"`  // 设置项类型，值为input/select/checkbox/switch等
	Options     string  `json:"options,omitempty" gorm:"type:text"`     // 可选值，仅在Type为select/checkbox时有效，格式为JSON数组
	Cols        int     `json:"cols,omitempty" gorm:"default:1"`        // 占用列数，默认1，最大12
}

const (
	SettingKeyTheme     = "theme"      // 主题，值为light/dark/auto
	SettingKeyAutoStart = "auto_start" // 是否开机自启，值为true或false
	SettingKeyLanguage  = "language"   // 语言，值为auto/zh-CN/en-US
)
