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

// UserLevelController struct
type UserLevel struct {
	base
}

// Index 用户等级 列表页
func (ulc *UserLevel) Index(c *gin.Context) {
	ulc.getBaseData(c)
	var userLevelService services.UserLevelService
	data, pagination := userLevelService.GetPaginateData(ulc.admin["per_page"].(int), ulc.gQueryParams)
	ulc.data["data"] = data
	ulc.data["paginate"] = pagination

	c.HTML(http.StatusOK, "user_level/index.html", ulc.data)
}

// Export 导出
func (ulc *UserLevel) Export(c *gin.Context) {
	ulc.getBaseData(c)
	exportData := c.PostForm("export_data")
	if exportData == "1" {
		var userLevelService services.UserLevelService
		data := userLevelService.GetExportData(ulc.gQueryParams)
		header := []string{"ID", "名称", "简介", "是否启用", "创建时间"}
		body := [][]string{}
		for _, item := range data {
			record := []string{
				strconv.Itoa(item.Id),
				item.Name,
				item.Description,
			}
			if item.Status == 1 {
				record = append(record, "是")
			} else {
				record = append(record, "否")
			}
			record = append(record, item.CreatedAt.Format(utils.TimeLayout))
			body = append(body, record)
		}
		c.Header("a", "b")
		exceloffice.ExportData(header, body, "user_level-"+time.Now().Format("2006-01-02-15-04-05"), "", "", c)
		return
	}

	response.Error(c)
}

// Add 用户等级-添加界面
func (ulc *UserLevel) Add(c *gin.Context) {
	ulc.getBaseData(c)
	c.HTML(http.StatusOK, "user_level/add.html", ulc.data)
}

// Create 用户等级-添加
func (ulc *UserLevel) Create(c *gin.Context) {
	ulc.getBaseData(c)
	var userLevelForm formvalidate.UserLevelForm
	if err := c.ShouldBind(&userLevelForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	// 自定义验证逻辑
	if err := validate.Struct(&userLevelForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	//处理图片上传
	_, err := c.FormFile("img")
	if err == nil {
		var attachmentService services.AttachmentService
		attachmentInfo, err := attachmentService.Upload(c, "img", ulc.loginUser.Id, 0)
		if err != nil || attachmentInfo == nil {
			response.ErrorWithMessage(err.Error(), c)
		} else {
			userLevelForm.Img = attachmentInfo.Url
		}
	}

	var userLevelService services.UserLevelService
	insertID := userLevelService.Create(&userLevelForm)

	url := global.URL_BACK
	if userLevelForm.IsCreate == 1 {
		url = global.URL_RELOAD
	}

	if insertID > 0 {
		response.SuccessWithMessageAndUrl("添加成功", url, c)
	} else {
		response.Error(c)
	}
}

// Edit 用户等级-修改界面
func (ulc *UserLevel) Edit(c *gin.Context) {
	ulc.getBaseData(c)
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		response.ErrorWithMessage("Param is error.", c)
		return
	}

	var userLevelService services.UserLevelService

	userLevel := userLevelService.GetUserLevelById(id)
	if userLevel == nil {
		response.ErrorWithMessage("Not Found Info By Id.", c)
		return
	}

	ulc.data["data"] = userLevel

	c.HTML(http.StatusOK, "user_level/edit.html", ulc.data)
}

// Update 用户等级-修改
func (ulc *UserLevel) Update(c *gin.Context) {
	var userLevelForm formvalidate.UserLevelForm
	if err := c.ShouldBind(&userLevelForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	if userLevelForm.Id <= 0 {
		response.ErrorWithMessage("Params is Error.", c)
		return
	}

	// 自定义验证逻辑
	if err := validate.Struct(&userLevelForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}
	_, err := c.FormFile("img")
	if err == nil {
		//处理图片上传
		var attachmentService services.AttachmentService
		attachmentInfo, err := attachmentService.Upload(c, "img", ulc.loginUser.Id, 0)
		if err != nil || attachmentInfo == nil {
			response.ErrorWithMessage(err.Error(), c)
			return
		} else {
			userLevelForm.Img = attachmentInfo.Url
		}
	}

	var userLevelService services.UserLevelService
	num := userLevelService.Update(&userLevelForm)

	if num > 0 {
		response.Success(c)
	} else {
		response.Error(c)
	}
}

// Enable 启用
func (ulc *UserLevel) Enable(c *gin.Context) {
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

	var userLevelService services.UserLevelService
	num := userLevelService.Enable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Disable 禁用
func (ulc *UserLevel) Disable(c *gin.Context) {
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

	var userLevelService services.UserLevelService
	num := userLevelService.Disable(idArr)
	if num > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}

// Del 删除
func (ulc *UserLevel) Del(c *gin.Context) {
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
		return
	}

	noDeletionID := new(models.UserLevel).NoDeletionId()

	m := go2.Intersect(noDeletionID, idArr)

	if len(noDeletionID) > 0 && len(m) > 0 {
		response.ErrorWithMessage("ID为"+strings.Join(utils.IntArrToStringArr(noDeletionID), ",")+"的数据无法删除!", c)
		return
	}

	var userLevelService services.UserLevelService
	count := userLevelService.Del(idArr)

	if count > 0 {
		response.SuccessWithMessageAndUrl("操作成功", global.URL_RELOAD, c)
	} else {
		response.Error(c)
	}
}
