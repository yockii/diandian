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
		Options:     `[{"title": "亮色", "value": "light"}, {"title": "暗色", "value": "dark"}, {"title": "自动", "value": "auto"}]`,
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
		Options:     `[{"title": "简体中文", "value": "zh-CN"}, {"title": "English", "value": "en-US"}, {"title": "自动", "value": "auto"}]`,
		Cols:        6,
	}).FirstOrCreate(&model.Setting{})
}
