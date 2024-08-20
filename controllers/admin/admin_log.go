package admin

import (
	"gin-admin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminLogController struct.
type AdminLog struct {
	base
}

// Index index.
func (alc *AdminLog) Index(c *gin.Context) {
	alc.getBaseData(c)
	var (
		adminLogService  services.AdminLogService
		adminUserService services.AdminUserService
	)
	data, pagination := adminLogService.GetPaginateData(alc.admin["per_page"].(int), alc.gQueryParams)

	alc.data["admin_user_list"] = adminUserService.GetAllAdminUser()

	alc.data["data"] = data
	alc.data["paginate"] = pagination
	c.HTML(http.StatusOK, "admin_log/index.html", alc.data)
}
