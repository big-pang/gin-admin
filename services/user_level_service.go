package services

import (
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/utils/page"
	"net/url"
	"time"
)

// UserLevelService struct
type UserLevelService struct {
	baseService
}

// GetPaginateData 通过分页获取user_level
func (uls *UserLevelService) GetPaginateData(listRows int, params url.Values) ([]*models.UserLevel, page.Pagination) {
	//搜索、查询字段赋值
	uls.SearchField = append(uls.SearchField, new(models.UserLevel).SearchField()...)

	var userLevel []*models.UserLevel
	o := global.DB.Model(new(models.UserLevel))
	err := uls.PaginateAndScopeWhere(o, listRows, params).Find(&userLevel).Error
	if err != nil {
		return nil, uls.Pagination
	}
	return userLevel, uls.Pagination
}

// GetExportData 获取导出数据
func (uls *UserLevelService) GetExportData(params url.Values) []*models.UserLevel {
	//搜索、查询字段赋值
	uls.SearchField = append(uls.SearchField, new(models.UserLevel).SearchField()...)
	var userLevel []*models.UserLevel
	o := global.DB.Model(new(models.UserLevel))
	err := uls.ScopeWhere(o, params).Find(&userLevel).Error
	if err != nil {
		return nil
	}
	return userLevel
}

// Create 新增用户等级
func (*UserLevelService) Create(form *formvalidate.UserLevelForm) int {
	userLevel := models.UserLevel{
		Name:        form.Name,
		Description: form.Description,
		Status:      int8(form.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if form.Img != "" {
		userLevel.Img = form.Img
	}
	err := global.DB.Create(&userLevel).Error

	if err == nil {
		return userLevel.Id
	}
	return 0
}

// Update 更新用户等级
func (*UserLevelService) Update(form *formvalidate.UserLevelForm) int {
	o := global.DB
	userLevel := models.UserLevel{Id: form.Id}
	if o.First(&userLevel).Error == nil {
		userLevel.Name = form.Name
		userLevel.Description = form.Description
		userLevel.Status = int8(form.Status)
		userLevel.UpdatedAt = time.Now()
		if form.Img != "" {
			userLevel.Img = form.Img
		}
		userLevel.Name = form.Name
		result := o.Updates(&userLevel)
		if result.Error == nil {
			return int(result.RowsAffected)
		}
		return 0
	}
	return 0
}

// GetUserLevelById 根据id获取一条user_level数据
func (*UserLevelService) GetUserLevelById(id int) *models.UserLevel {
	userLevel := models.UserLevel{Id: id}
	err := global.DB.First(&userLevel).Error
	if err != nil {
		return nil
	}
	return &userLevel
}

// GetUserLevel 获取所有用户等级
func (*UserLevelService) GetUserLevel() []*models.UserLevel {
	var userLevels []*models.UserLevel
	err := global.DB.Find(&userLevels).Error
	if err == nil {
		return userLevels
	}
	return nil
}

// Enable 启用
func (*UserLevelService) Enable(ids []int) int {
	result := global.DB.Model(new(models.UserLevel)).Where("id in ?", ids).Updates(models.Params{
		"status": 1,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Disable 禁用
func (*UserLevelService) Disable(ids []int) int {
	result := global.DB.Model(new(models.UserLevel)).Where("id in ?", ids).Updates(models.Params{
		"status": 0,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Del 删除
func (*UserLevelService) Del(ids []int) int {
	result := global.DB.Delete(&models.UserLevel{}, "id in ?", ids)
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}
