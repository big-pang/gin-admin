package services

import (
	"encoding/json"
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/utils/encrypter"
	"gin-admin/utils/page"
	"net/url"

	"github.com/gin-gonic/gin"
)

// AdminLogService struct
type AdminLogService struct {
	baseService
}

// LoginLog 登录日志
func (*AdminLogService) LoginLog(loginUserID int, ctx *gin.Context) {

}

// GetCount 获取admin_log 总数
func (*AdminLogService) GetCount() int {
	var count int64
	err := global.DB.Model(&models.AdminLog{}).Count(&count).Error
	if err != nil {
		global.LOG.Error(err.Error())
		return 0
	}
	return int(count)
}

// CreateAdminLog 创建操作日志
func (*AdminLogService) CreateAdminLog(loginUser *models.AdminUser, menu *models.AdminMenu, url string, c *gin.Context) {
	var adminLog models.AdminLog

	if loginUser == nil {
		adminLog.AdminUserId = 0
	} else {
		adminLog.AdminUserId = loginUser.Id
	}
	adminLog.Name = menu.Name
	adminLog.LogMethod = menu.LogMethod
	adminLog.Url = url
	adminLog.LogIp = c.ClientIP()

	//开启事务
	to := global.DB.Begin()

	err := to.Create(&adminLog).Error
	if err != nil {
		to.Rollback()
		global.LOG.Error(err.Error())
		return
	}
	//adminLogData数据表添加数据
	jsonData, _ := json.Marshal(c.Request.PostForm)
	cryptData := encrypter.Encrypt(jsonData, []byte(global.CONFIG.Base.LogAesKey))
	var adminLogData models.AdminLogData
	adminLogData.AdminLogId = int(adminLog.Id)
	adminLogData.Data = cryptData
	err = to.Create(&adminLogData).Error
	if err != nil {
		to.Rollback()
		global.LOG.Error(err.Error())
		return
	}
	to.Commit()
}

// GetPaginateData 获取所有adminuser
func (als *AdminLogService) GetPaginateData(listRows int, params url.Values) ([]*models.AdminLog, page.Pagination) {
	//搜索、查询字段赋值
	als.SearchField = append(als.SearchField, new(models.AdminLog).SearchField()...)
	als.WhereField = append(als.WhereField, new(models.AdminLog).WhereField()...)
	als.TimeField = append(als.TimeField, new(models.AdminLog).TimeField()...)

	var adminLog []*models.AdminLog
	o := global.DB.Model(new(models.AdminLog))
	err := als.PaginateAndScopeWhere(o, listRows, params).Find(&adminLog).Error
	if err != nil {
		return nil, als.Pagination
	}
	return adminLog, als.Pagination
}
