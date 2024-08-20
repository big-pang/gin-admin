package services

import (
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/utils/page"
	"net/url"
	"strings"
)

// AdminRoleService struct
type AdminRoleService struct {
	baseService
}

// GetCount 获取admin_role 总数
func (*AdminRoleService) GetCount() int {
	var count int64
	global.DB.Model(&models.AdminRole{}).Count(&count)
	return int(count)
}

// GetAllData 获取所有admin role
func (*AdminRoleService) GetAllData() []*models.AdminRole {
	var adminRoles []*models.AdminRole
	global.DB.Find(&adminRoles)
	return adminRoles
}

// GetPaginateData 分页获取adminrole
func (ars *AdminRoleService) GetPaginateData(listRows int, params url.Values) ([]*models.AdminRole, page.Pagination) {
	//搜索、查询字段赋值
	ars.SearchField = append(ars.SearchField, new(models.AdminRole).SearchField()...)

	var adminRole []*models.AdminRole
	o := global.DB.Model(&models.AdminRole{})
	err := ars.PaginateAndScopeWhere(o, listRows, params).Find(&adminRole).Error
	if err != nil {
		return nil, ars.Pagination
	}
	return adminRole, ars.Pagination
}

// IsExistName 名称验重
func (*AdminRoleService) IsExistName(name string, id int) bool {
	var count int64
	if id == 0 {
		global.DB.Model(&models.AdminRole{}).Where("name", name).Count(&count)
	} else {
		global.DB.Model(&models.AdminRole{}).Where("id != ?", id).Where("name", name).Count(&count)
	}
	return count > 0
}

// Create 创建角色
func (*AdminRoleService) Create(form *formvalidate.AdminRoleForm) int {
	adminRole := models.AdminRole{
		Name:        form.Name,
		Description: form.Description,
		Url:         "1,2,18",
		Status:      form.Status,
	}

	err := global.DB.Create(&adminRole).Error
	if err != nil {
		return 0
	}
	return adminRole.Id
}

// GetAdminRoleById 通过id获取菜单信息
func (*AdminRoleService) GetAdminRoleById(id int) *models.AdminRole {
	var adminRole models.AdminRole
	err := global.DB.Where("id = ?", id).First(&adminRole).Error

	if err == nil {
		return &adminRole
	}
	return nil
}

// Update 更新角色信息
func (*AdminRoleService) Update(form *formvalidate.AdminRoleForm) int {
	result := global.DB.Model(&models.AdminRole{}).Where("id", form.Id).Updates(models.Params{
		"name":        form.Name,
		"description": form.Description,
		"status":      form.Status,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Del 删除角色
func (*AdminRoleService) Del(ids []int) int {
	result := global.DB.Delete(&models.AdminRole{}, ids)
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Enable 启用角色
func (*AdminRoleService) Enable(ids []int) int {
	result := global.DB.Model(&models.AdminRole{}).Where("id in ?", ids).Updates(models.Params{
		"status": 1,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Disable 禁用角色
func (*AdminRoleService) Disable(ids []int) int {
	result := global.DB.Model(&models.AdminRole{}).Where("id in ?", ids).Updates(models.Params{
		"status": 0,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// StoreAccess 授权菜单
func (*AdminRoleService) StoreAccess(id int, url []string) int {
	result := global.DB.Model(new(models.AdminRole)).Where("id", id).Updates(models.Params{
		"url": strings.Join(url, ","),
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}
