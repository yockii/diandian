package service

import (
	"changeme/background/app"
	"changeme/background/database"
	"changeme/background/model"
	"changeme/background/util"
)

func InitializeData() {
	if !app.IsInitializeSuccess() {
		return
	}

	// 初始化设置类数据
	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyAutoStart,
	}).Attrs(&model.Setting{
		Value: util.StringPtr("false"),
	}).Assign(&model.Setting{
		GroupName:   "基础",
		Name:        "开机自启",
		Desc:        "是否随系统启动而启动",
		OrderNum:    1,
		Showable:    util.BoolPtr(false),
		SettingType: "switch",
		Cols:        6,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyTheme,
	}).Attrs(&model.Setting{
		Value: util.StringPtr("auto"),
	}).Assign(&model.Setting{
		GroupName:   "基础",
		Name:        "主题",
		Desc:        "应用主题，自动为跟随系统",
		OrderNum:    2,
		Showable:    util.BoolPtr(true),
		SettingType: "select",
		Options:     `[{"label": "亮色", "value": "light"}, {"label": "暗色", "value": "dark"}, {"label": "自动", "value": "auto"}]`,
		Cols:        6,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyLanguage,
	}).Attrs(&model.Setting{
		Value: util.StringPtr("auto"),
	}).Assign(&model.Setting{
		GroupName:   "基础",
		Name:        "语言",
		Desc:        "应用语言，自动为跟随系统",
		OrderNum:    3,
		Showable:    util.BoolPtr(false),
		SettingType: "select",
		Options:     `[{"label": "简体中文", "value": "zh-CN"}, {"label": "English", "value": "en-US"}, {"label": "自动", "value": "auto"}]`,
		Cols:        6,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyLlmBaseUrl,
	}).Assign(&model.Setting{
		GroupName:   "大模型",
		Name:        "Base URL",
		Desc:        "访问大模型的基础URL",
		OrderNum:    1,
		Showable:    util.BoolPtr(true),
		SettingType: "input",
		Cols:        24,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyLlmToken,
	}).Assign(&model.Setting{
		GroupName:   "大模型",
		Name:        "Token",
		Desc:        "访问大模型API的Token",
		OrderNum:    2,
		Showable:    util.BoolPtr(true),
		SettingType: "password",
		Cols:        24,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyTextModel,
	}).Assign(&model.Setting{
		GroupName:   "大模型",
		Name:        "文本模型",
		Desc:        "要使用的文本生成模型(为空则使用视觉模型完成任务)",
		OrderNum:    3,
		Showable:    util.BoolPtr(true),
		SettingType: "input",
		Cols:        12,
	}).FirstOrCreate(&model.Setting{})

	database.DB.Model(&model.Setting{}).Where(&model.Setting{
		Key: model.SettingKeyVlModel,
	}).Assign(&model.Setting{
		GroupName:   "大模型",
		Name:        "视觉模型",
		Desc:        "要使用的视觉生成模型",
		OrderNum:    4,
		Showable:    util.BoolPtr(true),
		SettingType: "input",
		Cols:        12,
	}).FirstOrCreate(&model.Setting{})
}
