package admin

import (
	"gin-admin/formvalidate"
	"gin-admin/formvalidate/validate"
	"gin-admin/global"
	"gin-admin/global/response"
	"gin-admin/services"
	"gin-admin/utils"
	"strconv"

	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

var adminLogService services.AdminLogService

type Auth struct{}

func (ac *Auth) Login(c *gin.Context) {
	//获取登录配置信息
	loginConfig := struct {
		Token      string
		Captcha    string
		Background string
	}{
		Token:      global.CONFIG.Login.Token,
		Captcha:    global.CONFIG.Login.Captcha,
		Background: global.CONFIG.Login.Background,
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
		"admin": gin.H{
			"name":  global.CONFIG.Base.Name,
			"title": "登录",
		},
		"login_config":   loginConfig,
		"captcha":        utils.GetCaptcha(),
		"debug":          global.CONFIG.Base.Debug,
		"csrf_token":     csrf.Token(c.Request),
		csrf.TemplateTag: csrf.TemplateField(c.Request),
	})

}

func (ac *Auth) GetCaptcha(c *gin.Context) {
	captchaId := c.Param("captchaId")
	c.Header("Content-Type", "image/png")
	captcha.WriteImage(c.Writer, captchaId, captcha.StdWidth, captcha.StdHeight)
}

func (ac *Auth) RefreshCaptcha(c *gin.Context) {
	captchaID := c.Param("captchaId")
	res := map[string]interface{}{
		"isNew": false,
	}
	if captchaID == "" {
		res["msg"] = "参数错误."
	}

	isReload := captcha.Reload(captchaID)
	if isReload {
		res["captchaId"] = captchaID
	} else {
		res["isNew"] = true
		res["captcha"] = utils.GetCaptcha()
	}

	c.JSON(http.StatusOK, res)
}

func (ac *Auth) CheckLogin(c *gin.Context) {

	// 数据校验
	var loginForm formvalidate.LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	// 看是否需要校验验证码
	isCaptcha, _ := strconv.Atoi(global.CONFIG.Login.Captcha)
	if isCaptcha > 0 {
		if loginForm.Captcha == "" {
			response.ErrorWithMessage("请输入验证码.", c)
			return
		}

		if ok := captcha.VerifyString(loginForm.CaptchaId, loginForm.Captcha); !ok {
			response.ErrorWithMessage("验证码错误.", c)
		}
	}

	// 自定义验证逻辑
	if err := validate.Struct(&loginForm); err != nil {
		response.ErrorWithMessage("验证失败: "+err.Error(), c)
		return
	}

	// 基础验证通过后，进行用户验证
	var adminUserService services.AdminUserService
	loginUser, err := adminUserService.CheckLogin(loginForm, c)
	if err != nil {
		response.ErrorWithMessage(err.Error(), c)
		return
	}

	// 登录日志记录
	adminLogService.LoginLog(loginUser.Id, c)

	redirect, _ := c.Get("redirect")
	if redirectStr, ok := redirect.(string); ok && redirectStr != "" {
		response.SuccessWithMessageAndUrl("登录成功", redirectStr, c)
	} else {
		response.SuccessWithMessageAndUrl("登录成功", "/admin/index/index", c)
	}
}

// Logout 退出登录
func (ac *Auth) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(global.LOGIN_USER)
	session.Save()
	c.SetCookie(global.LOGIN_USER_ID, "", -1, "/", "", false, false)
	c.SetCookie(global.LOGIN_USER_ID_SIGN, "", -1, "/", "", false, false)
	response.SuccessWithMessageAndUrl("退出成功", "/admin/auth/login", c)
}
