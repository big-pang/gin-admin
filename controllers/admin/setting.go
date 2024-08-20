package admin

import (
	"encoding/json"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Setting struct
type Setting struct {
	base
}

// Admin 设置
func (sc *Setting) Admin(c *gin.Context) {
	sc.getBaseData(c)
	var settingService services.SettingService
	data := settingService.Show(1)

	sc.data["data_config"] = data
	c.HTML(http.StatusOK, "setting/show.html", sc.data)

}

// show 展示单个配置信息
func (sc *Setting) show(id int, c *gin.Context) {
	sc.getBaseData(c)
	var settingService services.SettingService
	data := settingService.Show(id)

	sc.data["data_config"] = data
	c.HTML(http.StatusOK, "setting/show.html", sc.data)
}

// Update 设置中心-更新设置
func (sc *Setting) Update(c *gin.Context) {
	sc.getBaseData(c)
	id := c.PostForm("id")

	if id == "" {
		response.ErrorWithMessage("参数错误.", c)
		return
	}

	var settingService services.SettingService
	idInt, _ := strconv.Atoi(id)
	setting := settingService.GetSettingInfoById(idInt)

	if setting == nil {
		response.ErrorWithMessage("无法更新配置信息", c)
		return
	}

	err := json.Unmarshal([]byte(setting.Content), &setting.ContentStrut)
	if err != nil {
		response.ErrorWithMessage("无法更新配置信息", c)
		return
	}

	for key, value := range setting.ContentStrut {
		switch value.Type {
		case "image", "file":
			//单个文件上传
			var attachmentService services.AttachmentService
			attachmentInfo, err := attachmentService.Upload(c, value.Field, sc.loginUser.Id, 0)
			if err == nil && attachmentInfo != nil {
				//图片上传成功
				setting.ContentStrut[key].Content = attachmentInfo.Url
			}
		case "multi_file", "multi_image":
			//多个文件上传
			var attachmentService services.AttachmentService
			attachments, err := attachmentService.UploadMulti(c, value.Field, sc.loginUser.Id, 0)
			if err == nil && attachments != nil {
				var urls []string
				for _, atta := range attachments {
					urls = append(urls, atta.Url)
				}
				if len(urls) > 0 {
					urlsByte, err := json.Marshal(&urls)
					if err == nil {
						setting.ContentStrut[key].Content = string(urlsByte)
					}
				}
			}
		default:
			setting.ContentStrut[key].Content = c.PostForm(value.Field)
		}
	}

	//修改内容
	contentStrutByte, err := json.Marshal(&setting.ContentStrut)
	if err == nil {
		affectRow := settingService.UpdateSettingInfoToContent(idInt, string(contentStrutByte))
		if affectRow > 0 {
			//更新全局配置
			settingService.LoadOrUpdateGlobalBaseConfig(setting)
			response.SuccessWithMessageAndUrl("修改成功", global.URL_RELOAD, c)
			return
		} else {
			response.ErrorWithMessage("没有可更新的信息", c)
			return
		}
	} else {
		response.ErrorWithMessage("修改失败 err:"+err.Error(), c)
		return
	}

}
