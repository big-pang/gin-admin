package models

import (
	"crypto/sha1"

	"fmt"
	"gin-admin/global"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminUser struct
type AdminUser struct {
	Id        int            `gorm:"column(id);auto;size(11);description(表ID)" json:"id"`
	Username  string         `gorm:"column(username);size(30);description(用户名)" json:"username"`
	Password  string         `gorm:"column(password);size(255);description(密码)" json:"password"`
	Nickname  string         `gorm:"column(nickname);size(30);description(昵称)" json:"nickname"`
	Avatar    string         `gorm:"column(avatar);size(255);description(头像)" json:"avatar"`
	Role      string         `gorm:"column(role);size(200);description(角色)" json:"role"`
	Status    int8           `gorm:"column(status);size(1);description(是否启用 0：否 1：是)" json:"status"`
	DeletedAt gorm.DeletedAt `gorm:"column(deleted_at);type(timestamp);default(NULL);description(删除时间)" json:"deleted_at"`
}

// TableName 自定义table 名称
func (*AdminUser) TableName() string {
	return "admin_user"
}

// SearchField 定义模型的可搜索字段
func (*AdminUser) SearchField() []string {
	return []string{"nickname", "username"}
}

// NoDeletionId 禁止删除的数据id
func (*AdminUser) NoDeletionId() []int {
	return []int{1}
}

// WhereField 定义模型可作为条件的字段
func (*AdminUser) WhereField() []string {
	return []string{}
}

// TimeField 定义可做为时间范围查询的字段
func (*AdminUser) TimeField() []string {
	return []string{}
}

// 在init中注册定义的model
func init() {

}

// GetSignStrByAdminUser 获取加密字符串，用在登录的时候加密处理
// func (adminUser *AdminUser) GetSignStrByAdminUser(ctx *context.Context) string {
// 	ua := ctx.Input.Header("user-agent")
// 	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%d%s%s", adminUser.Id, adminUser.Username, ua))))
// }

// GetAuthUrl 获取已授权url
func (adminUser *AdminUser) GetAuthUrl() map[string]int {
	var (
		urlArr     []string
		authURLArr []string
	)
	authURL := make(map[string]int)
	err := global.DB.Model(&AdminRole{}).Select("url").Where("id IN ? and status = 1", strings.Split(adminUser.Role, ",")).Pluck("url", &urlArr).Error
	if err == nil {
		urlIDStr := ""
		for k, row := range urlArr {
			if k == 0 {
				urlIDStr = row
			} else {
				urlIDStr += "," + row
			}
		}

		urlIDArr := strings.Split(urlIDStr, ",")

		if len(urlIDStr) > 0 {

			err := global.DB.Model(AdminMenu{}).Select("url").Where("id IN ?", urlIDArr).Pluck("url", &authURLArr).Error
			if err == nil {
				for k, row := range authURLArr {
					authURL[row] = k
				}
			}
		}
		return authURL
	}

	return authURL
}

// GetSignStrByAdminUser 获取加密字符串，用在登录的时候加密处理
func (adminUser *AdminUser) GetSignStrByAdminUser(c *gin.Context) string {
	ua := c.Request.Header.Get("User-Agent")
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%d%s%s", adminUser.Id, adminUser.Username, ua))))
}

// GetShowMenu 获取当前用户已授权的显示菜单
func (adminUser *AdminUser) GetShowMenu() map[int]Params {
	var maps []map[string]interface{}
	returnMaps := make(map[int]Params)
	if adminUser.Id == 1 {
		err := global.DB.Model(&AdminMenu{}).Where("is_show", 1).Order("sort_id, id").Find(&maps).Error
		if err == nil {

			for _, m := range maps {
				returnMaps[m["id"].(int)] = m
			}
			return returnMaps
		}
		return map[int]Params{}
	}

	var list []string
	err := global.DB.Model(&AdminRole{}).Select("url").Where("id IN ? and status = 1", strings.Split(adminUser.Role, ",")).Pluck("url", &list).Error
	if err == nil {
		var urlIDArr []string
		for _, m := range list {
			urlIDArr = append(urlIDArr, strings.Split(m, ",")...)
		}
		err := global.DB.Model(AdminMenu{}).Select("id", "parent_id", "name", "url", "icon", "sort_id").Where("id IN ?", list).Where("is_show", 1).Order("sort_id, id").Find(&maps).Error

		if err == nil {
			for _, m := range maps {
				returnMaps[m["id"].(int)] = m
			}
			return returnMaps
		}
		return map[int]Params{}
	}
	return map[int]Params{}

}

// GetRoleText 用户角色名称
func (adminUser *AdminUser) GetRoleText() map[int]*AdminRole {
	roleIDArr := strings.Split(adminUser.Role, ",")
	var adminRole []*AdminRole
	err := global.DB.Select("id", "name").Where("id IN ?", roleIDArr).Find(&adminRole).Error

	if err != nil {
		return nil
	}
	adminRoleMap := make(map[int]*AdminRole)
	for _, v := range adminRole {
		adminRoleMap[v.Id] = v
	}
	return adminRoleMap
}
