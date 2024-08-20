package admin

import (
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
	"github.com/gin-gonic/gin"
)

// AdminMenuController struct.
type AdminMenu struct {
	base
}

// Index 菜单首页.
func (amc *AdminMenu) Index(c *gin.Context) {
	amc.getBaseData(c)
	var adminTreeService services.AdminTreeService
	amc.data["data"] = adminTreeService.AdminMenuTree()

	c.HTML(http.StatusOK, "admin_menu/index.html", amc.data)
}

// Add 添加菜单界面.
func (amc *AdminMenu) Add(c *gin.Context) {
	amc.getBaseData(c)
	var adminTreeService services.AdminTreeService
	parentID, _ := strconv.Atoi(c.Query("parent_id"))
	parents := adminTreeService.Menu(parentID, 0)

	amc.data["parents"] = parents
	amc.data["log_method"] = new(models.AdminMenu).GetLogMethod()

	c.HTML(http.StatusOK, "admin_menu/add.html", amc.data)
}

// Create 添加菜单.
func (amc *AdminMenu) Create(c *gin.Context) {
	var adminMenuService services.AdminMenuService
	adminMenuForm := formvalidate.AdminMenuForm{}

	if err := c.ShouldBind(&adminMenuForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	//去除Url前后两侧的空格
	if adminMenuForm.Url != "" {
		adminMenuForm.Url = strings.TrimSpace(adminMenuForm.Url)
	}

	// 自定义验证逻辑
	if err := validate.Struct(&adminMenuForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//添加之前url验重
	if adminMenuService.IsExistUrl(adminMenuForm.Url, adminMenuForm.Id) {
		response.ErrorWithMessage("url【"+adminMenuForm.Url+"】已经存在.", c)
		return
	}

	//创建
	_, err := adminMenuService.Create(&adminMenuForm)
	if err != nil {
		response.Error(c)
		return
	}

	url := global.URL_BACK
	if adminMenuForm.IsCreate == 1 {
		url = global.URL_RELOAD
	}

	response.SuccessWithMessageAndUrl("添加成功", url, c)
}

// Edit 编辑菜单界面.
func (amc *AdminMenu) Edit(c *gin.Context) {
	amc.getBaseData(c)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var (
		adminMenuService services.AdminMenuService
		adminTreeService services.AdminTreeService
	)

	adminMenu := adminMenuService.GetAdminMenuById(id)
	if adminMenu == nil {
		response.ErrorWithMessage("Not Found Info By Id.", c)
		return
	}

	parentID := adminMenu.ParentId
	parents := adminTreeService.Menu(parentID, 0)

	amc.data["data"] = adminMenu
	amc.data["parents"] = parents
	amc.data["log_method"] = new(models.AdminMenu).GetLogMethod()

	c.HTML(http.StatusOK, "admin_menu/edit.html", amc.data)
}

// Update 菜单更新.
func (amc *AdminMenu) Update(c *gin.Context) {
	var adminMenuService services.AdminMenuService
	adminMenuForm := formvalidate.AdminMenuForm{}

	if err := c.ShouldBind(&adminMenuForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	//去除Url前后两侧的空格
	if adminMenuForm.Url != "" {
		adminMenuForm.Url = strings.TrimSpace(adminMenuForm.Url)
	}

	// 自定义验证逻辑
	if err := validate.Struct(&adminMenuForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//添加之前url验重
	if adminMenuService.IsExistUrl(adminMenuForm.Url, adminMenuForm.Id) {
		response.ErrorWithMessage("url【"+adminMenuForm.Url+"】已经存在.", c)
		return
	}

	count := adminMenuService.Update(&adminMenuForm)

	if count > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}
}

// Del 删除.
func (amc *AdminMenu) Del(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	ids := make([]int, 0)
	var idArr []int
	var formData IdsFormData
	if id <= 0 {
		c.Bind(&formData)
		ids = formData.Ids
	} else {
		idArr = append(idArr, id)
	}

	if len(ids) > 0 {
		idArr = ids
	}

	if len(idArr) == 0 {
		response.ErrorWithMessage("参数id错误.", c)
		return
	}

	var adminMenuService services.AdminMenuService
	//判断是否有子菜单
	if adminMenuService.IsChildMenu(idArr) {
		response.ErrorWithMessage("有子菜单不可删除！", c)
		return
	}

	noDeletionID := new(models.AdminMenu).NoDeletionId()

	m := go2.Intersect(noDeletionID, idArr)

	if len(noDeletionID) > 0 && len(m) > 0 {
		response.ErrorWithMessage("ID为"+strings.Join(utils.IntArrToStringArr(noDeletionID), ",")+"的数据无法删除!", c)
		return
	}

	count := adminMenuService.Del(idArr)

	if count > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}
