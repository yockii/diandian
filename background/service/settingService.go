package service

import (
	"changeme/background/app"
	"changeme/background/database"
	"changeme/background/model"
)

type SettingService struct{}

func (s *SettingService) AllSettings() (list []*model.Setting, err error) {
	err = database.DB.Order("group_name, order_num").Find(&list).Error
	return
}

func (s *SettingService) SaveSetting(setting *model.Setting) (err error) {
	err = database.DB.Save(setting).Error
	if err != nil {
		return
	}

	// 获取该设置
	var currentSetting *model.Setting
	err = database.DB.Where("id = ?", setting.ID).First(&currentSetting).Error
	if err != nil {
		return
	}
	// 触发相应事件
	switch currentSetting.Key {
	case model.SettingKeyTheme:
		app.EmitEvent("theme-change", currentSetting.Value)
	}
	return
}

func (s *SettingService) GetThemeSetting() (*model.Setting, error) {
	var setting model.Setting
	err := database.DB.Where("key = ?", model.SettingKeyTheme).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}
