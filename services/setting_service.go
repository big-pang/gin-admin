package services

import (
	"encoding/json"
	"gin-admin/global"
	"gin-admin/models"
)

// SettingService struct
type SettingService struct {
	baseService
}

// Show 显示设置
func (settingService *SettingService) Show(id int) []*models.Setting {
	data := settingService.getDataBySettingGroupId(id)

	var settingFormService SettingFormService

	for key, value := range data {
		//contentNew := ""
		//value.Content转为json
		var contents []*models.Content
		if value.Content == "" {
			continue
		}
		err := json.Unmarshal([]byte(value.Content), &contents)

		if err != nil {
			continue
		}

		var contentNew []*models.Content
		for _, content := range contents {
			content.Form = settingFormService.GetFieldForm(content.Type, content.Name, content.Field, content.Content, content.Option)
			contentNew = append(contentNew, content)
		}
		data[key].ContentStrut = contentNew
	}

	return data
}

// getDataBySettingGroupId 根据设置分组id获取多个设置信息
func (*SettingService) getDataBySettingGroupId(settingGroupId int) []*models.Setting {
	var settings []*models.Setting
	err := global.DB.Where("setting_group_id", settingGroupId).Find(&settings).Error
	if err != nil {
		return nil
	}
	return settings
}

// GetSettingInfoById 根据设置id，获取对应的setting info
func (*SettingService) GetSettingInfoById(id int) *models.Setting {
	setting := models.Setting{Id: id}
	global.DB.First(&setting)
	return &setting
}

// UpdateSettingInfoToContent 根据id修改content的内容
func (*SettingService) UpdateSettingInfoToContent(id int, content string) int {
	result := global.DB.Where("id", id).Updates(models.Setting{
		Content: content,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// LoadOrUpdateGlobalBaseConfig 加载或者更新全局登录、系统配置信息
func (*SettingService) LoadOrUpdateGlobalBaseConfig(setting *models.Setting) bool {
	if setting == nil {
		return false
	}

	if setting.Code == "base" {
		for _, content := range setting.ContentStrut {
			switch content.Field {
			case "name":
				global.CONFIG.Base.Name = content.Content
			case "short_name":
				global.CONFIG.Base.ShortName = content.Content
			case "author":
				global.CONFIG.Base.Author = content.Content
			case "version":
				global.CONFIG.Base.Version = content.Content
			case "link":
				global.CONFIG.Base.Link = content.Content
			}
		}
	} else if setting.Code == "login" {
		for _, content := range setting.ContentStrut {
			switch content.Field {
			case "token":
				global.CONFIG.Login.Token = content.Content
			case "captcha":
				global.CONFIG.Login.Captcha = content.Content
			case "background":
				global.CONFIG.Login.Background = content.Content
			}
		}
	} else if setting.Code == "index" {
		for _, content := range setting.ContentStrut {
			switch content.Field {
			case "password_warning":
				global.CONFIG.Base.PasswordWarning = content.Content
			case "show_notice":
				global.CONFIG.Base.ShowNotice = content.Content
			case "notice_content":
				global.CONFIG.Base.NoticeContent = content.Content
			}
		}
	} else {
		return false
	}

	return true
}
