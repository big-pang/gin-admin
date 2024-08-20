package admin

import (
	"gin-admin/formvalidate"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/models"
	"gin-admin/services"
	"gin-admin/utils"
	"net/http"
	"strconv"
	"strings"

	"gin-admin/formvalidate/validate"

	go2 "github.com/adam-hanna/arrayOperations"
	"github.com/gin-gonic/gin"
)

// AdminRoleController struct.
type AdminRole struct {
	base
}

// Index 角色管理首页.
func (arc *AdminRole) Index(c *gin.Context) {
	arc.getBaseData(c)
	var adminRoleService services.AdminRoleService
	data, pagination := adminRoleService.GetPaginateData(arc.admin["per_page"].(int), arc.gQueryParams)

	arc.data["data"] = data
	arc.data["paginate"] = pagination

	c.HTML(http.StatusOK, "admin_role/index.html", arc.data)
}

// Add 角色管理-添加界面.
func (arc *AdminRole) Add(c *gin.Context) {
	arc.getBaseData(c)
	c.HTML(http.StatusOK, "admin_role/add.html", arc.data)
}

// Create 角色管理-添加角色.
func (arc *AdminRole) Create(c *gin.Context) {
	var adminRoleForm formvalidate.AdminRoleForm
	if err := c.ShouldBind(&adminRoleForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
	}

	// 自定义验证逻辑
	if err := validate.Struct(&adminRoleForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	var adminRoleService services.AdminRoleService

	//名称验重
	if adminRoleService.IsExistName(adminRoleForm.Name, 0) {
		response.ErrorWithMessage("名称已存在.", c)
	}

	url := global.URL_BACK
	if adminRoleForm.IsCreate == 1 {
		url = global.URL_RELOAD
	}

	//添加
	insertID := adminRoleService.Create(&adminRoleForm)
	if insertID > 0 {
		response.SuccessWithMessageAndUrl("添加成功", url, c)
	} else {
		response.Error(c)
	}
}

// Edit 菜单管理-角色管理-修改界面.
func (arc *AdminRole) Edit(c *gin.Context) {
	arc.getBaseData(c)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var (
		adminRoleService services.AdminRoleService
	)

	adminRole := adminRoleService.GetAdminRoleById(id)
	if adminRole == nil {
		response.ErrorWithMessage("Not Found Info By Id.", c)
		return
	}

	arc.data["data"] = adminRole

	c.HTML(http.StatusOK, "admin_role/edit.html", arc.data)
}

// Update 菜单管理-角色管理-修改.
func (arc *AdminRole) Update(c *gin.Context) {
	var adminRoleForm formvalidate.AdminRoleForm
	if err := c.ShouldBind(&adminRoleForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
	}

	//id验证
	if adminRoleForm.Id == 0 {
		response.ErrorWithMessage("id错误.", c)
		return
	}

	//字段验证
	if err := validate.Struct(adminRoleForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	var adminRoleService services.AdminRoleService

	//名称验重
	if adminRoleService.IsExistName(adminRoleForm.Name, adminRoleForm.Id) {
		response.ErrorWithMessage("名称已存在.", c)
	}

	//修改
	num := adminRoleService.Update(&adminRoleForm)
	if num > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}
}

// Del 删除.
func (arc *AdminRole) Del(c *gin.Context) {
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

	var adminRoleService services.AdminRoleService

	noDeletionID := new(models.AdminRole).NoDeletionId()

	m := go2.Intersect(noDeletionID, idArr)

	if len(noDeletionID) > 0 && len(m) > 0 {
		response.ErrorWithMessage("ID为"+strings.Join(utils.IntArrToStringArr(noDeletionID), ",")+"的数据无法删除!", c)
	}

	count := adminRoleService.Del(idArr)

	if count > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Enable 启用.
func (arc *AdminRole) Enable(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	ids := make([]int, 0)
	var formData IdsFormData
	var idArr []int
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
		response.ErrorWithMessage("请选择启用的角色.", c)
		return
	}

	var adminRoleService services.AdminRoleService
	num := adminRoleService.Enable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Disable 禁用.
func (arc *AdminRole) Disable(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	ids := make([]int, 0)
	var idArr []int
	var formData IdsFormData
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
		response.ErrorWithMessage("请选择禁用的角色.", c)
		return
	}

	var adminRoleService services.AdminRoleService
	num := adminRoleService.Disable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Access 菜单管理-角色管理-角色授权界面.
func (arc *AdminRole) Access(c *gin.Context) {
	arc.getBaseData(c)
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var (
		adminRoleService services.AdminRoleService
		adminMenuService services.AdminMenuService
		adminTreeService services.AdminTreeService
	)

	data := adminRoleService.GetAdminRoleById(id)
	menu := adminMenuService.AllMenu()

	menuMap := make(map[int]models.Params)

	for _, adminMenu := range menu {
		id := adminMenu.Id
		if menuMap[id] == nil {
			menuMap[id] = make(models.Params)
		}
		menuMap[id]["Id"] = id
		menuMap[id]["ParentId"] = adminMenu.ParentId
		menuMap[id]["Name"] = adminMenu.Name
		menuMap[id]["Url"] = adminMenu.Url
		menuMap[id]["Icon"] = adminMenu.Icon
		menuMap[id]["IsShow"] = adminMenu.IsShow
		menuMap[id]["SortId"] = adminMenu.SortId
		menuMap[id]["LogMethod"] = adminMenu.LogMethod
	}

	html := adminTreeService.AuthorizeHtml(menuMap, strings.Split(data.Url, ","))

	arc.data["data"] = data
	arc.data["html"] = html

	c.HTML(http.StatusOK, "admin_role/access.html", arc.data)
}

// AccessOperate 菜单管理-角色管理-角色授权.
func (arc *AdminRole) AccessOperate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		response.ErrorWithMessage("Params is Error.", c)
	}

	url := make([]string, 0)
	url = c.PostFormArray("url[]")
	if len(url) == 0 {
		response.ErrorWithMessage("请选择授权的菜单", c)
	}

	if !utils.InArrayForString(url, "1") {
		response.ErrorWithMessage("首页权限必选", c)
	}

	var adminRoleService services.AdminRoleService
	num := adminRoleService.StoreAccess(id, url)
	if num > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}

}
