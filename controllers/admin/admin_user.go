package admin

import (
	"encoding/base64"
	"gin-admin/formvalidate"
	"gin-admin/formvalidate/validate"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/models"
	"gin-admin/services"
	"gin-admin/utils"
	"net/http"
	"strconv"
	"strings"

	go2 "github.com/adam-hanna/arrayOperations"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

type AdminUser struct {
	base
}

// Index 用户管理-首页
func (auc *AdminUser) Index(c *gin.Context) {
	auc.getBaseData(c)
	var adminUserService services.AdminUserService
	adminUsers, pagination := adminUserService.GetPaginateData(auc.admin["per_page"].(int), auc.gQueryParams)
	auc.data["data"] = adminUsers
	auc.data["paginate"] = pagination
	c.HTML(200, "admin_user/index.html", auc.data)
}

// Add 用户管理-添加界面
func (auc *AdminUser) Add(c *gin.Context) {
	auc.getBaseData(c)
	var adminRoleService services.AdminRoleService
	roles := adminRoleService.GetAllData()
	auc.data["roles"] = roles
	c.HTML(200, "admin_user/add.html", auc.data)
}

// Create 用户管理-添加界面
func (auc *AdminUser) Create(c *gin.Context) {
	var adminUserForm formvalidate.AdminUserForm
	if err := c.ShouldBind(&adminUserForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
	}
	// 自定义验证逻辑
	if err := validate.Struct(&adminUserForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//账号验重
	var adminUserService services.AdminUserService
	if adminUserService.IsExistName(strings.TrimSpace(adminUserForm.Username), 0) {
		response.ErrorWithMessage("账号已经存在", c)
	}
	//默认头像
	adminUserForm.Avatar = "/static/admin/images/avatar.png"

	insertID := adminUserService.Create(&adminUserForm)
	url := global.URL_BACK
	if adminUserForm.IsCreate == 1 {
		url = global.URL_RELOAD
	}
	if insertID > 0 {
		response.SuccessWithMessageAndUrl("添加成功", url, c)
	} else {
		response.Error(c)
	}
}

// Edit 系统管理-用户管理-修改界面
func (auc *AdminUser) Edit(c *gin.Context) {
	auc.getBaseData(c)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var (
		adminUserService services.AdminUserService
		adminRoleService services.AdminRoleService
	)
	adminUser := adminUserService.GetAdminUserById(id)
	if adminUser == nil {
		response.ErrorWithMessage("Not Found Info By Id.", c)
	}

	roles := adminRoleService.GetAllData()
	auc.data["roles"] = roles
	auc.data["data"] = adminUser
	auc.data["role_arr"] = strings.Split(adminUser.Role, ",")

	c.HTML(http.StatusOK, "admin_user/edit.html", auc.data)

}

// Update 系统管理-用户管理-修改
func (auc *AdminUser) Update(c *gin.Context) {
	var adminUserForm formvalidate.AdminUserForm
	if err := c.ShouldBind(&adminUserForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
	}
	// 自定义验证逻辑
	if err := validate.Struct(&adminUserForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//账号验重
	var adminUserService services.AdminUserService
	if adminUserService.IsExistName(strings.TrimSpace(adminUserForm.Username), adminUserForm.Id) {
		response.ErrorWithMessage("账号已经存在", c)
	}

	numb := adminUserService.Update(&adminUserForm)

	if numb > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}
}

// Enable 启用
func (auc *AdminUser) Enable(c *gin.Context) {
	idStr := c.Query("id")
	ids := make([]int, 0)

	var formData IdsFormData
	var idArr []int
	if idStr == "" {
		c.Bind(&formData)
		ids = formData.Ids
	} else {
		id, _ := strconv.Atoi(idStr)
		idArr = append(idArr, id)
	}

	if len(ids) > 0 {
		idArr = ids
	}

	if len(idArr) == 0 {
		response.ErrorWithMessage("请选择启用的用户.", c)
		return
	}

	var adminUserService services.AdminUserService
	num := adminUserService.Enable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}

}

// Disable 禁用
func (auc *AdminUser) Disable(c *gin.Context) {
	idStr := c.Query("id")
	ids := make([]int, 0)

	var formData IdsFormData
	var idArr []int
	if idStr == "" {
		c.Bind(&formData)
		ids = formData.Ids
	} else {
		id, _ := strconv.Atoi(idStr)
		idArr = append(idArr, id)
	}

	if len(ids) > 0 {
		idArr = ids
	}

	if len(idArr) == 0 {
		response.ErrorWithMessage("请选择禁用的用户.", c)
		return
	}

	var adminUserService services.AdminUserService
	num := adminUserService.Disable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Del 删除
func (auc *AdminUser) Del(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	ids := make([]int, 0)

	var formData IdsFormData
	var idArr []int
	if id == 0 {
		c.Bind(&formData)
		ids = formData.Ids
	} else {
		idArr = append(idArr, id)
	}

	if len(ids) > 0 {
		idArr = ids
	}

	if len(idArr) == 0 {
		response.ErrorWithMessage("请选择禁用的用户.", c)
		return
	}

	noDeletionID := new(models.AdminUser).NoDeletionId()

	m := go2.Intersect(noDeletionID, idArr)

	if len(noDeletionID) > 0 && len(m) > 0 {
		response.ErrorWithMessage("ID为"+strings.Join(utils.IntArrToStringArr(noDeletionID), ",")+"的数据无法删除!", c)
		return
	}

	var adminUserService services.AdminUserService
	count := adminUserService.Del(idArr)
	if count > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Profile 系统管理-个人资料
func (auc *AdminUser) Profile(c *gin.Context) {
	auc.getBaseData(c)
	c.HTML(http.StatusOK, "admin_user/profile.html", auc.data)
}

// UpdateNickName 系统管理-个人资料-修改昵称
func (auc *AdminUser) UpdateNickName(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}
	nickname := strings.TrimSpace(c.GetString("nickname"))

	if nickname == "" {
		response.ErrorWithMessage("参数错误", c)
	}

	var adminUserService services.AdminUserService
	num := adminUserService.UpdateNickName(id, nickname)

	if num > 0 {
		//修改成功后，更新session的登录用户信息
		loginAdminUser := adminUserService.GetAdminUserById(id)
		session := sessions.Default(c)
		session.Set(global.LOGIN_USER, loginAdminUser)
		err := session.Save()
		if err != nil {
			response.ErrorWithMessage("session保存失败:"+err.Error(), c)
		} else {
			response.SuccessWithMessageAndUrl("修改成功", global.URL_RELOAD, c)
		}
	} else {
		response.Error(c)
	}
}

// UpdatePassword 系统管理-个人资料-修改密码
func (auc *AdminUser) UpdatePassword(c *gin.Context) {
	auc.getBaseData(c)
	type UpdateStruct struct {
		Id            int    `form:"id" binding:"required"`
		Password      string `form:"password" binding:"required"`
		NewPassword   string `form:"new_password" binding:"required"`
		ReNewPassword string `form:"renew_password" binding:"required"`
	}
	updateDate := UpdateStruct{}
	err := c.ShouldBind(&updateDate)

	if err != nil || updateDate.Id == 0 || updateDate.Password == "" || updateDate.NewPassword == "" || updateDate.ReNewPassword == "" {
		response.ErrorWithMessage("Bad Parameter.", c)
	}

	if updateDate.NewPassword != updateDate.ReNewPassword {
		response.ErrorWithMessage("两次输入的密码不一致.", c)
	}

	if updateDate.Password == updateDate.NewPassword {
		response.ErrorWithMessage("新密码与旧密码一致，无需修改", c)
	}

	loginUserPassword, err := base64.StdEncoding.DecodeString(auc.loginUser.Password)

	if err != nil {
		response.ErrorWithMessage("err:"+err.Error(), c)
	}

	if !utils.PasswordVerify(updateDate.Password, string(loginUserPassword)) {
		response.ErrorWithMessage("当前密码不正确", c)
	}

	var adminUserService services.AdminUserService
	num := adminUserService.UpdatePassword(updateDate.Id, updateDate.NewPassword)
	if num > 0 {
		response.SuccessWithMessageAndUrl("修改成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// UpdateAvatar 系统管理-个人资料-修改头像
func (auc *AdminUser) UpdateAvatar(c *gin.Context) {
	_, err := c.FormFile("avatar")
	if err != nil {
		response.ErrorWithMessage("上传头像错误"+err.Error(), c)
	}

	var (
		attachmentService services.AttachmentService
		adminUserService  services.AdminUserService
	)
	attachmentInfo, err := attachmentService.Upload(c, "avatar", auc.loginUser.Id, 0)
	if err != nil || attachmentInfo == nil {
		response.ErrorWithMessage(err.Error(), c)
	} else {
		//头像上传成功，更新用户的avatar头像信息
		num := adminUserService.UpdateAvatar(auc.loginUser.Id, attachmentInfo.Url)
		if num > 0 {
			//修改成功后，更新session的登录用户信息
			loginAdminUser := adminUserService.GetAdminUserById(auc.loginUser.Id)
			session := sessions.Default(c)
			session.Set(global.LOGIN_USER, *loginAdminUser)
			session.Save()
			err := session.Save()
			if err != nil {
				response.ErrorWithMessage("session保存失败:"+err.Error(), c)
			} else {
				response.SuccessWithMessageAndUrl("修改成功", global.URL_RELOAD, c)
			}

		} else {
			response.Error(c)
		}
	}
}
