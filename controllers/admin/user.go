package admin

import (
	"gin-admin/formvalidate"
	"gin-admin/formvalidate/validate"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/models"
	"gin-admin/services"
	"gin-admin/utils"
	"gin-admin/utils/exceloffice"
	"net/http"
	"strconv"
	"strings"
	"time"

	go2 "github.com/adam-hanna/arrayOperations"
	"github.com/gin-gonic/gin"
)

// UserController struct
type User struct {
	base
}

// Index 用户等级 列表页
func (uc *User) Index(c *gin.Context) {
	uc.getBaseData(c)
	var userService services.UserService
	var userLevelService services.UserLevelService

	//获取用户等级
	userLevel := userLevelService.GetUserLevel()
	userLevelMap := make(map[int]string)
	for _, item := range userLevel {
		userLevelMap[item.Id] = item.Name
	}

	data, pagination := userService.GetPaginateData(uc.admin["per_page"].(int), uc.gQueryParams)

	uc.data["data"] = data
	uc.data["paginate"] = pagination
	uc.data["user_level_map"] = userLevelMap

	c.HTML(http.StatusOK, "user/index.html", uc.data)
}

// Export 导出
func (uc *User) Export(c *gin.Context) {
	uc.getBaseData(c)
	exportData := c.PostForm("export_data")
	if exportData == "1" {
		var userService services.UserService
		var userLevelService services.UserLevelService
		userLevel := userLevelService.GetUserLevel()
		userLevelMap := make(map[int]string)
		for _, item := range userLevel {
			userLevelMap[item.Id] = item.Name
		}

		data := userService.GetExportData(uc.gQueryParams)
		header := []string{"ID", "头像", "用户等级", "用户名", "手机号", "昵称", "是否启用", "创建时间"}
		body := [][]string{}
		for _, item := range data {
			record := []string{
				strconv.Itoa(item.Id),
				item.Avatar,
			}
			userLevelName, ok := userLevelMap[item.UserLevelId]
			if ok {
				record = append(record, userLevelName)
			}
			record = append(record, item.Username)
			record = append(record, item.Mobile)
			record = append(record, item.Nickname)

			if item.Status == 1 {
				record = append(record, "是")
			} else {
				record = append(record, "否")
			}
			record = append(record, item.CreatedAt.Format(utils.TimeLayout))
			body = append(body, record)
		}
		c.Header("a", "b")
		exceloffice.ExportData(header, body, "user-"+time.Now().Format("2006-01-02-15-04-05"), "", "", c)
		return
	}

	response.Error(c)
}

// Add 用户-添加界面
func (uc *User) Add(c *gin.Context) {
	uc.getBaseData(c)
	var userLevelService services.UserLevelService

	//获取用户等级
	userLevel := userLevelService.GetUserLevel()

	uc.data["user_level_list"] = userLevel
	c.HTML(http.StatusOK, "user/add.html", uc.data)
}

// Create 添加用户
func (uc *User) Create(c *gin.Context) {
	var userForm formvalidate.UserForm
	if err := c.ShouldBind(&userForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	// 自定义验证逻辑
	if err := validate.Struct(&userForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//处理图片上传
	_, err := c.FormFile("avatar")
	if err == nil {
		var attachmentService services.AttachmentService
		attachmentInfo, err := attachmentService.Upload(c, "avatar", uc.loginUser.Id, 0)
		if err != nil || attachmentInfo == nil {
			response.ErrorWithMessage(err.Error(), c)
			return
		} else {
			userForm.Avatar = attachmentInfo.Url
		}
	}

	var userService services.UserService
	insertID := userService.Create(&userForm)

	url := global.URL_BACK
	if userForm.IsCreate == 1 {
		url = global.URL_RELOAD
	}

	if insertID > 0 {
		response.SuccessWithMessageAndUrl("添加成功", url, c)
	} else {
		response.Error(c)
	}
}

// Edit 用户-修改界面
func (uc *User) Edit(c *gin.Context) {
	uc.getBaseData(c)
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var userService services.UserService

	user := userService.GetUserById(id)
	if user == nil {
		response.ErrorWithMessage("Not Found Info By Id.", c)
		return
	}

	//获取用户等级
	var userLevelService services.UserLevelService
	userLevel := userLevelService.GetUserLevel()

	uc.data["user_level_list"] = userLevel
	uc.data["data"] = user

	c.HTML(http.StatusOK, "user/edit.html", uc.data)
}

// Update 用户-修改
func (uc *User) Update(c *gin.Context) {
	var userForm formvalidate.UserForm
	if err := c.ShouldBind(&userForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	if userForm.Id <= 0 {
		response.ErrorWithMessage("Params is Error.", c)
		return
	}

	// 自定义验证逻辑
	if err := validate.Struct(&userForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	_, err := c.FormFile("avatar")
	if err == nil {
		//处理图片上传
		var attachmentService services.AttachmentService
		attachmentInfo, err := attachmentService.Upload(c, "avatar", uc.loginUser.Id, 0)
		if err != nil || attachmentInfo == nil {
			response.ErrorWithMessage(err.Error(), c)
			return
		} else {
			userForm.Avatar = attachmentInfo.Url
		}
	}

	var userService services.UserService
	num := userService.Update(&userForm)

	if num > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}
}

// Enable 启用
func (uc *User) Enable(c *gin.Context) {
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
		response.ErrorWithMessage("请选择用户等级.", c)
		return
	}

	var userService services.UserService
	num := userService.Enable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
		return
	} else {
		response.Error(c)
	}
}

// Disable 禁用
func (uc *User) Disable(c *gin.Context) {
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

	var userService services.UserService
	num := userService.Disable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Del 删除
func (uc *User) Del(c *gin.Context) {
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
		response.ErrorWithMessage("参数id错误.", c)
	}

	noDeletionID := new(models.User).NoDeletionId()

	m := go2.Intersect(noDeletionID, idArr)

	if len(noDeletionID) > 0 && len(m) > 0 {
		response.ErrorWithMessage("ID为"+strings.Join(utils.IntArrToStringArr(noDeletionID), ",")+"的数据无法删除!", c)
		return
	}

	var userService services.UserService
	count := userService.Del(idArr)

	if count > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}
