package services

import (
	"encoding/base64"
	"errors"
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/utils"
	"gin-admin/utils/page"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminUserService struct
type AdminUserService struct {
	baseService
}

// GetAdminUserById 根据id获取一条admin_user数据
func (*AdminUserService) GetAdminUserById(id int) *models.AdminUser {
	adminUser := &models.AdminUser{Id: id}
	err := global.DB.Find(adminUser).Error
	if err != nil {
		return nil
	}
	return adminUser
}

// AuthCheck 权限检测
func (*AdminUserService) AuthCheck(url string, authExcept map[string]int, loginUser *models.AdminUser) bool {
	authURL := loginUser.GetAuthUrl()
	if utils.KeyInMap(url, authExcept) || utils.KeyInMap(url, authURL) {
		return true
	}
	return false
}

// CheckLogin 用户登录验证
func (*AdminUserService) CheckLogin(loginForm formvalidate.LoginForm, c *gin.Context) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := global.DB.Where("username = ?", loginForm.Username).First(&adminUser).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	decodePasswdStr, err := base64.StdEncoding.DecodeString(adminUser.Password)

	if err != nil || !utils.PasswordVerify(loginForm.Password, string(decodePasswdStr)) {
		return nil, errors.New("密码错误")
	}

	if adminUser.Status != 1 {
		return nil, errors.New("用户被冻结")
	}
	session := sessions.Default(c)
	// ctx.Output.Session(global.LOGIN_USER, adminUser)
	session.Set(global.LOGIN_USER, adminUser)
	err = session.Save()
	if err != nil {
		global.LOG.Sugar().Error(err)
	}

	if loginForm.Remember != "" {
		c.SetCookie(global.LOGIN_USER_ID, strconv.Itoa(adminUser.Id), 7200, "/", "", false, false)
		c.SetCookie(global.LOGIN_USER_ID_SIGN, adminUser.GetSignStrByAdminUser(c), 7200, "/", "", false, false)
	} else {
		c.SetCookie(global.LOGIN_USER_ID, strconv.Itoa(adminUser.Id), -1, "/", "", false, false)
		c.SetCookie(global.LOGIN_USER_ID_SIGN, adminUser.GetSignStrByAdminUser(c), -1, "/", "", false, false)
	}

	return &adminUser, nil

}

// GetCount 获取admin_user 总数
func (*AdminUserService) GetCount() int {
	var count int64
	global.DB.Model(&models.AdminUser{}).Count(&count)
	return int(count)
}

// GetAllAdminUser 获取所有adminuser
func (*AdminUserService) GetAllAdminUser() []*models.AdminUser {
	var adminUser []*models.AdminUser

	err := global.DB.Find(&adminUser).Error
	if err != nil {
		return nil
	}
	return adminUser
}

// UpdateNickName 系统管理-个人资料-修改昵称
func (*AdminUserService) UpdateNickName(id int, nickname string) int {
	result := global.DB.Model(&models.AdminUser{Id: id}).Updates(models.AdminUser{Nickname: nickname})
	if result.Error != nil || result.RowsAffected <= 0 {
		return 0
	}
	return int(result.RowsAffected)
}

// UpdatePassword 修改密码
func (*AdminUserService) UpdatePassword(id int, newPassword string) int {
	newPasswordForHash, err := utils.PasswordHash(newPassword)

	if err != nil {
		return 0
	}
	result := global.DB.Model(&models.AdminUser{Id: id}).Update("password", base64.StdEncoding.EncodeToString([]byte(newPasswordForHash)))

	if result.Error != nil || result.RowsAffected <= 0 {
		return 0
	}
	return int(result.RowsAffected)
}

// UpdateAvatar 系统管理-个人资料-修改头像
func (*AdminUserService) UpdateAvatar(id int, avatar string) int {
	result := global.DB.Model(&models.AdminUser{Id: id}).Update("avatar", avatar)
	if result.Error != nil || result.RowsAffected <= 0 {
		return 0
	}
	return int(result.RowsAffected)
}

// GetPaginateData 通过分页获取adminuser
func (aus *AdminUserService) GetPaginateData(listRows int, params url.Values) ([]*models.AdminUser, page.Pagination) {
	//搜索、查询字段赋值
	aus.SearchField = append(aus.SearchField, new(models.AdminUser).SearchField()...)

	var adminUser []*models.AdminUser
	o := global.DB.Model(new(models.AdminUser))
	err := aus.PaginateAndScopeWhere(o, listRows, params).Find(&adminUser).Error
	if err != nil {
		return nil, aus.Pagination
	}
	return adminUser, aus.Pagination
}

// IsExistName 名称验重
func (*AdminUserService) IsExistName(username string, id int) bool {
	var count int64
	if id == 0 {
		global.DB.Model(new(models.AdminUser)).Where("username", username).Count(&count)
		return count > 0
	}
	global.DB.Model(new(models.AdminUser)).Where("id != ?", id).Where("username", username).Count(&count)
	return count > 0
}

// Create 新增admin user用户
func (*AdminUserService) Create(form *formvalidate.AdminUserForm) int {
	newPasswordForHash, err := utils.PasswordHash(form.Password)
	if err != nil {
		return 0
	}

	adminUser := models.AdminUser{
		Username: form.Username,
		Password: base64.StdEncoding.EncodeToString([]byte(newPasswordForHash)),
		Nickname: form.Nickname,
		Avatar:   form.Avatar,
		Role:     strings.Join(form.Roles, ","),
		Status:   int8(form.Status),
	}
	result := global.DB.Create(&adminUser)

	if result.Error == nil {
		return adminUser.Id
	}
	return 0
}

// Update 更新用户信息
func (*AdminUserService) Update(form *formvalidate.AdminUserForm) int {
	o := global.DB
	adminUser := models.AdminUser{Id: form.Id}
	if o.First(&adminUser).RowsAffected > 0 {
		adminUser.Username = form.Username
		adminUser.Nickname = form.Nickname
		adminUser.Role = strings.Join(form.Roles, ",")
		adminUser.Status = int8(form.Status)
		if adminUser.Password != form.Password {
			newPasswordForHash, err := utils.PasswordHash(form.Password)
			if err == nil {
				adminUser.Password = base64.StdEncoding.EncodeToString([]byte(newPasswordForHash))
			}
		}
		result := o.Updates(&adminUser)
		if result.Error == nil {
			return int(result.RowsAffected)
		}
		return 0
	}
	return 0
}

// Enable 启用用户
func (*AdminUserService) Enable(ids []int) int {
	result := global.DB.Model(&models.AdminUser{}).Where("id in (?)", ids).Update("status", 1)

	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Disable 禁用用户
func (*AdminUserService) Disable(ids []int) int {
	result := global.DB.Model(&models.AdminUser{}).Where("id in (?)", ids).Update("status", 0)

	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Del 删除用户
func (*AdminUserService) Del(ids []int) int {
	result := global.DB.Delete(&models.AdminUser{}, ids)

	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}
