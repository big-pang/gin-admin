package middleware

import (
	"fmt"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/models"
	"gin-admin/services"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddle 中间件
func AuthMiddle() gin.HandlerFunc {

	//不需要验证的url
	authExcept := map[string]int{
		"admin/auth/login":           0,
		"admin/auth/check_login":     1,
		"admin/auth/logout":          2,
		"admin/auth/captcha":         3,
		"admin/editor/server":        4,
		"admin/auth/refresh_captcha": 5,
	}

	//登录认证中间件过滤器
	return func(c *gin.Context) {
		url := strings.TrimLeft(c.Request.URL.Path, "/")

		//需要进行登录验证
		if !isAuthExceptUrl(strings.ToLower(url), authExcept) {
			//验证是否登录
			loginUser, isLogin := isLogin(c)
			if !isLogin {
				response.ErrorWithMessageAndUrl("未登录", "/admin/auth/login", c)
				c.Abort()
				return
			}

			//验证，是否有权限访问
			var adminUserService services.AdminUserService
			if loginUser.Id != 1 && !adminUserService.AuthCheck(url, authExcept, loginUser) {
				errorBackURL := global.URL_CURRENT
				if c.Request.Method == "GET" {
					errorBackURL = ""
				}
				response.ErrorWithMessageAndUrl("无权限", errorBackURL, c)
				c.Abort()

			}
		}

		checkAuth, _ := strconv.Atoi(c.PostForm("check_auth"))

		if checkAuth == 1 {
			response.Success(c)
			return
		}

	}
}

// 判断是否是不需要验证登录的url,只针对admin模块路由的判断
func isAuthExceptUrl(url string, m map[string]int) bool {
	urlArr := strings.Split(url, "/")
	if len(urlArr) > 3 {
		url = fmt.Sprintf("%s/%s/%s", urlArr[0], urlArr[1], urlArr[2])
	}
	_, ok := m[url]
	if ok {
		return true
	}
	return false
}

// 是否登录
func isLogin(c *gin.Context) (*models.AdminUser, bool) {
	loginUser, ok := sessions.Default(c).Get(global.LOGIN_USER).(models.AdminUser)
	if !ok {
		loginUserIDStr, _ := c.Cookie(global.LOGIN_USER_ID)
		loginUserIDSign, _ := c.Cookie(global.LOGIN_USER_ID_SIGN)

		if loginUserIDStr != "" && loginUserIDSign != "" {
			loginUserID, _ := strconv.Atoi(loginUserIDStr)
			var adminUserService services.AdminUserService
			loginUserPointer := adminUserService.GetAdminUserById(loginUserID)

			if loginUserPointer != nil && loginUserPointer.GetSignStrByAdminUser(c) == loginUserIDSign {
				session := sessions.Default(c)
				session.Set(global.LOGIN_USER, *loginUserPointer)
				err := session.Save()
				if err != nil {
					response.ErrorWithMessage("登录失败", c)
					c.Abort()
				}
				return loginUserPointer, true
			}
		}
		return nil, false
	}

	return &loginUser, true
}
