package middleware

import (
	"gin-admin/global"
	"gin-admin/models"
	"gin-admin/services"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// baseController struct
type baseController struct {
}

var (
	//后台变量
	admin map[string]interface{}
	//当前用户
	loginUser models.AdminUser
	//参数
	gQueryParams url.Values
)

// Prepare 父控制器初始化
func Prepare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 访问url
		requestURL := strings.ToLower(strings.TrimLeft(c.Request.URL.Path, "/"))
		data := map[string]any{}
		// query参数
		// 只有分页首页列表时才会使用
		if c.Request.Method == "GET" {
			gQueryParams, _ = url.ParseQuery(c.Request.URL.RawQuery)
			gQueryParams.Set("queryParamUrl", c.Request.URL.String())
			if len(gQueryParams) > 0 {
				for k, val := range gQueryParams {
					v, ok := strconv.Atoi(val[0])
					if ok == nil {
						data[k] = v
					} else {
						data[k] = val[0]
					}
				}
			}
		}

		// 登录用户
		var isOk bool
		loginUser, isOk = sessions.Default(c).Get(global.LOGIN_USER).(models.AdminUser)
		// 基础变量
		runMode := global.CONFIG.Base.Mode
		if runMode == "development" {
			data["debug"] = true
		} else {
			data["debug"] = false
		}
		data["cookie_prefix"] = ""
		perPageStr, _ := c.Request.Cookie("admin_per_page")
		// 每页预览的数量
		var perPage int
		if perPageStr == nil {
			perPage = 10
		} else {
			perPage, _ = strconv.Atoi(perPageStr.Value)
		}
		if perPage >= 100 {
			perPage = 100
		}

		// 记录日志
		var adminMenuService services.AdminMenuService
		adminMenu := adminMenuService.GetAdminMenuByUrl(requestURL)
		title := ""
		if adminMenu != nil {
			title = adminMenu.Name
			if strings.ToLower(adminMenu.LogMethod) == strings.ToLower(c.Request.Method) {
				var adminLogService services.AdminLogService
				adminLogService.CreateAdminLog(&loginUser, adminMenu, requestURL, c)
			}
		}

		// 左侧菜单
		menu := ""
		if "admin/auth/login" != requestURL && !(c.GetHeader("X-PJAX") == "true") && isOk {
			var adminTreeService services.AdminTreeService
			menu = adminTreeService.GetLeftMenu(requestURL, loginUser)
		}

		admin = map[string]interface{}{
			"pjax":            c.GetHeader("X-PJAX") == "true",
			"user":            &loginUser,
			"menu":            menu,
			"name":            global.CONFIG.Base.Name,
			"author":          global.CONFIG.Base.Author,
			"version":         global.CONFIG.Base.Version,
			"short_name":      global.CONFIG.Base.ShortName,
			"link":            global.CONFIG.Base.Link,
			"per_page":        perPage,
			"per_page_config": []int{10, 20, 30, 50, 100},
			"title":           title,
		}
		data["admin"] = admin
		data["gQueryParams"] = gQueryParams
		c.Set("gQueryParams", gQueryParams)
		c.Set("adminBaseData", data)
		c.Set("loginUser", loginUser)

	}

}
