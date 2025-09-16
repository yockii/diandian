package service

import (
	"diandian/background/app"
	"diandian/background/constant"
	"diandian/background/database"
	"diandian/background/model"
)

type SettingService struct{}

func (s *SettingService) AllSettings() (list []*model.Setting, err error) {
	err = database.DB.Order("group_name, order_num").Find(&list).Error
	return
}

func (s *SettingService) SaveSetting(setting *model.Setting) (err error) {
	err = database.DB.Select("value").Save(setting).Error
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
		app.EmitEvent(constant.EventThemeChanged, currentSetting.Value)
	case model.SettingKeyLlmBaseUrl, model.SettingKeyLlmToken, model.SettingKeyVlModel, model.SettingKeyTextModel:
		canWork, _ := s.CanWork()
		app.EmitEvent(constant.EventCanWorkChanged, canWork)
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

func (s *SettingService) CanWork() (bool, error) {
	var settings []*model.Setting
	err := database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmBaseUrl,
		model.SettingKeyLlmToken,
		model.SettingKeyVlModel,
		model.SettingKeyTextModel,
	}).Find(&settings).Error

	if err != nil {
		return false, err
	}

	if len(settings) < 4 {
		return false, nil
	}

	for _, setting := range settings {
		if setting.Value == nil || *setting.Value == "" {
			return false, nil
		}
	}

	return true, nil
}
