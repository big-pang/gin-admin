package admin

import (
	"gin-admin/models"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

type IdsFormData struct {
	Ids []int `form:"id[]"`
}

type base struct {
	data          gin.H
	loginUser     *models.AdminUser
	adminBaseData map[string]any
	gQueryParams  url.Values
	admin         map[string]any
}

func (bc *base) getBaseData(c *gin.Context) {
	bc.data = gin.H{}
	bc.loginUser, bc.adminBaseData, bc.gQueryParams = getAdminData(c)
	bc.data["login_user"] = bc.loginUser
	bc.admin = bc.adminBaseData["admin"].(map[string]any)

	bc.data["csrf_token"] = csrf.Token(c.Request)
	bc.data[csrf.TemplateTag] = csrf.TemplateField(c.Request)

	for k, v := range bc.adminBaseData {
		bc.data[k] = v
	}
}

func getAdminData(c *gin.Context) (*models.AdminUser, map[string]any, url.Values) {
	loginUser := models.AdminUser{}
	adminBaseData := make(map[string]any)
	gQueryParams := url.Values{}
	login_user, ok := c.Get("loginUser")
	if ok {
		loginUser = login_user.(models.AdminUser)
	}
	admin_base_data, ok := c.Get("adminBaseData")
	if ok {
		adminBaseData = admin_base_data.(map[string]any)
	}
	g_query_params, ok := c.Get("gQueryParams")
	if ok {
		gQueryParams = g_query_params.(url.Values)
	}
	return &loginUser, adminBaseData, gQueryParams
}
