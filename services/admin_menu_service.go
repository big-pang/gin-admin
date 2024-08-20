package services

import (
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/models"
)

// AdminMenuService struct
type AdminMenuService struct {
}

// GetAdminMenuByUrl 根据url获取admin_menu数据
func (*AdminMenuService) GetAdminMenuByUrl(url string) *models.AdminMenu {
	var adminMenu models.AdminMenu
	err := global.DB.Where("url =?", url).First(&adminMenu).Error
	if err == nil {
		return &adminMenu
	}
	return nil
}

// GetCount 获取admin_menu 总数
func (*AdminMenuService) GetCount() int {
	var count int64
	err := global.DB.Model(&models.AdminMenu{}).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

// AllMenu 获取所有菜单
func (*AdminMenuService) AllMenu() []*models.AdminMenu {
	var adminMenus []*models.AdminMenu
	err := global.DB.Model(&models.AdminMenu{}).Order("sort_id, id").Find(&adminMenus).Error
	if err == nil {
		return adminMenus
	}
	return nil
}

// Menu 除去当前id之外的所有菜单id
func (*AdminMenuService) Menu(currentID int) []map[string]interface{} {
	var adminMenusMap []map[string]interface{}
	global.DB.Where("id != ?", currentID).Order("sort_id, id").Find(&adminMenusMap)
	return adminMenusMap
}

// Create 创建菜单
func (*AdminMenuService) Create(form *formvalidate.AdminMenuForm) (int64, error) {
	adminMenu := models.AdminMenu{
		ParentId:  form.ParentId,
		Name:      form.Name,
		Url:       form.Url,
		Icon:      form.Icon,
		IsShow:    form.IsShow,
		SortId:    form.SortId,
		LogMethod: form.LogMethod,
	}
	err := global.DB.Create(&adminMenu).Error
	if err != nil {
		return 0, err
	}
	return int64(adminMenu.Id), nil
}

// Update 更新菜单
func (*AdminMenuService) Update(form *formvalidate.AdminMenuForm) int {
	result := global.DB.Where("id = ?", form.Id).Updates(models.AdminMenu{
		ParentId:  form.ParentId,
		Name:      form.Name,
		Url:       form.Url,
		Icon:      form.Icon,
		IsShow:    form.IsShow,
		SortId:    form.SortId,
		LogMethod: form.LogMethod,
	})

	if err := result.Error; err == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// IsExistUrl Url验重
func (*AdminMenuService) IsExistUrl(url string, id int) bool {
	var count int64
	if id == 0 {
		global.DB.Where("url = ?", url).Count(&count)
	} else {
		global.DB.Where("url = ?", url).Where("id != ?", id).Count(&count)
	}
	if count == 0 {
		return false
	}
	return true

}

// IsChildMenu 判断是否有子菜单
func (*AdminMenuService) IsChildMenu(ids []int) bool {
	var count int64
	global.DB.Where("parent_id in ?", ids).Count(&count)
	if count == 0 {
		return false
	} else {
		return true
	}

}

// Del 删除菜单
func (*AdminMenuService) Del(ids []int) int {
	result := global.DB.Where("id in ?", ids).Delete(new(models.AdminMenu))

	if err := result.Error; err == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// GetAdminMenuById 通过id获取菜单信息
func (*AdminMenuService) GetAdminMenuById(id int) *models.AdminMenu {
	var adminMenu models.AdminMenu
	err := global.DB.Where("id = ?", id).First(&adminMenu).Error
	if err == nil {
		return &adminMenu
	}
	return nil
}
