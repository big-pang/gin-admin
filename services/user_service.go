package services

import (
	"encoding/base64"
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/utils"
	"gin-admin/utils/page"
	"net/url"
	"time"
)

// UserService struct
type UserService struct {
	baseService
}

// GetPaginateData 通过分页获取user
func (us *UserService) GetPaginateData(listRows int, params url.Values) ([]*models.User, page.Pagination) {
	//搜索、查询字段赋值
	us.SearchField = append(us.SearchField, new(models.User).SearchField()...)

	var users []*models.User
	o := global.DB.Model(new(models.User))
	err := us.PaginateAndScopeWhere(o, listRows, params).Find(&users).Error
	if err != nil {
		return nil, us.Pagination
	}
	return users, us.Pagination
}

// Create 新增用户
func (*UserService) Create(form *formvalidate.UserForm) int {
	user := models.User{
		Username:    form.Username,
		Nickname:    form.Nickname,
		UserLevelId: form.UserLevelId,
		Mobile:      form.Mobile,
		Description: form.Description,
		Status:      int8(form.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if form.Avatar != "" {
		user.Avatar = form.Avatar
	}

	//密码加密
	newPasswordForHash, err := utils.PasswordHash(form.Password)
	if err != nil {
		return 0
	}
	user.Password = base64.StdEncoding.EncodeToString([]byte(newPasswordForHash))

	err = global.DB.Create(&user).Error

	if err == nil {
		return user.Id
	}
	return 0
}

// GetUserById 根据id获取一条user数据
func (*UserService) GetUserById(id int) *models.User {
	user := models.User{Id: id}
	err := global.DB.First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// Update 更新用户
func (*UserService) Update(form *formvalidate.UserForm) int {
	user := models.User{Id: form.Id}
	if global.DB.First(&user) == nil {

		//判断密码是否相等
		if user.Password != form.Password {
			newPasswordForHash, err := utils.PasswordHash(form.Password)
			if err == nil {
				user.Password = base64.StdEncoding.EncodeToString([]byte(newPasswordForHash))
			}
		}

		user.Username = form.Username
		user.Nickname = form.Nickname
		user.UserLevelId = form.UserLevelId
		user.Mobile = form.Mobile
		user.Description = form.Description
		user.Status = int8(form.Status)
		user.UpdatedAt = time.Now()

		if form.Avatar != "" {
			user.Avatar = form.Avatar
		}
		result := global.DB.Updates(&user)

		if result.Error == nil {
			return int(result.RowsAffected)
		}
		return 0
	}
	return 0
}

// Enable 启用
func (*UserService) Enable(ids []int) int {
	result := global.DB.Model(new(models.User)).Where("id in ?", ids).Updates(models.Params{
		"status": 1,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Disable 禁用
func (*UserService) Disable(ids []int) int {
	result := global.DB.Model(new(models.User)).Where("id in ?", ids).Updates(models.Params{
		"status": 0,
	})
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// Del 删除
func (*UserService) Del(ids []int) int {
	result := global.DB.Model(new(models.User)).Delete(&models.User{}, "id in ?", ids)
	if result.Error == nil {
		return int(result.RowsAffected)
	}
	return 0
}

// GetExportData 获取导出数据
func (us *UserService) GetExportData(params url.Values) []*models.User {
	//搜索、查询字段赋值
	us.SearchField = append(us.SearchField, new(models.User).SearchField()...)
	var user []*models.User
	o := global.DB.Model(new(models.User))
	err := us.ScopeWhere(o, params).Find(&user).Error
	if err != nil {
		return nil
	}
	return user
}
