package router

import (
	"fmt"
	"gin-admin/controllers/admin"
	"gin-admin/global"
	"gin-admin/middleware"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func adminRouter() {
	router.Static("/static", "static")
	// router.LoadHTMLGlob("views/admin/**/*")
	// 加载模板文件
	// 加载所有模板文件
	global.LOG.Debug("set HTMLRender .....")
	router.HTMLRender = loadTemplates("views/admin")

	initSession()
	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(middleware.CSRF())
		adminRouter.Use(middleware.SCRFToken())
		adminRouter.Use(middleware.Prepare())
		adminRouter.Use(middleware.AuthMiddle())

		adminAuthRouter := adminRouter.Group("/auth")
		{
			adminAuthRouter.GET("/login", admin.AdminGroupApp.Auth.Login)
			adminAuthRouter.GET("/logout", admin.AdminGroupApp.Auth.Logout)
			adminAuthRouter.POST("/check_login", admin.AdminGroupApp.Auth.CheckLogin)
			adminAuthRouter.GET("/captcha/:captchaId", admin.AdminGroupApp.Auth.GetCaptcha)
			adminAuthRouter.POST("/refresh_captcha", admin.AdminGroupApp.Auth.RefreshCaptcha)

		}

		//middleware.Prepare()

		//首页
		adminRouter.GET("/index/index", admin.AdminGroupApp.Index.Index)

		//用户管理
		adminRouter.GET("/admin_user/index", admin.AdminGroupApp.AdminUser.Index)
		//系统管理-个人资料
		adminRouter.GET("/admin_user/profile", admin.AdminGroupApp.AdminUser.Profile)
		//系统管理-个人资料-修改昵称
		adminRouter.POST("/admin_user/update_nickname", admin.AdminGroupApp.AdminUser.UpdateNickName)
		//系统管理-个人资料-修改密码
		adminRouter.POST("/admin_user/update_password", admin.AdminGroupApp.AdminUser.UpdatePassword)
		//系统管理-个人资料-修改头像
		adminRouter.POST("/admin_user/update_avatar", admin.AdminGroupApp.AdminUser.UpdateAvatar)
		//系统管理-用户管理-添加界面
		adminRouter.GET("/admin_user/add", admin.AdminGroupApp.AdminUser.Add)
		//系统管理-用户管理-添加
		adminRouter.POST("/admin_user/create", admin.AdminGroupApp.AdminUser.Create)
		//系统管理-用户管理-修改界面
		adminRouter.GET("/admin_user/edit", admin.AdminGroupApp.AdminUser.Edit)
		//系统管理-用户管理-修改
		adminRouter.POST("/admin_user/update", admin.AdminGroupApp.AdminUser.Update)
		//系统管理-用户管理-启用
		adminRouter.POST("/admin_user/enable", admin.AdminGroupApp.AdminUser.Enable)
		//系统管理-用户管理-禁用
		adminRouter.POST("/admin_user/disable", admin.AdminGroupApp.AdminUser.Disable)
		//系统管理-用户管理-删除
		adminRouter.POST("/admin_user/del", admin.AdminGroupApp.AdminUser.Del)

		//系统管理-角色管理
		adminRouter.GET("/admin_role/index", admin.AdminGroupApp.AdminRole.Index)
		//系统管理-角色管理-添加界面
		adminRouter.GET("/admin_role/add", admin.AdminGroupApp.AdminRole.Add)
		//系统管理-角色管理-添加
		adminRouter.POST("/admin_role/create", admin.AdminGroupApp.AdminRole.Create)
		//菜单管理-角色管理-修改界面
		adminRouter.GET("/admin_role/edit", admin.AdminGroupApp.AdminRole.Edit)
		//菜单管理-角色管理-修改
		adminRouter.POST("/admin_role/update", admin.AdminGroupApp.AdminRole.Update)
		//菜单管理-角色管理-删除
		adminRouter.POST("/admin_role/del", admin.AdminGroupApp.AdminRole.Del)
		//菜单管理-角色管理-启用角色
		adminRouter.POST("/admin_role/enable", admin.AdminGroupApp.AdminRole.Enable)
		//菜单管理-角色管理-禁用角色
		adminRouter.POST("/admin_role/disable", admin.AdminGroupApp.AdminRole.Disable)
		//菜单管理-角色管理-角色授权界面
		adminRouter.GET("/admin_role/access", admin.AdminGroupApp.AdminRole.Access)
		//菜单管理-角色管理-角色授权
		adminRouter.POST("/admin_role/access_operate", admin.AdminGroupApp.AdminRole.AccessOperate)

		//菜单管理
		adminRouter.GET("/admin_menu/index", admin.AdminGroupApp.AdminMenu.Index)
		//菜单管理-添加菜单-界面
		adminRouter.GET("/admin_menu/add", admin.AdminGroupApp.AdminMenu.Add)
		//菜单管理-添加菜单-创建
		adminRouter.POST("/admin_menu/create", admin.AdminGroupApp.AdminMenu.Create)
		//菜单管理-修改菜单-界面
		adminRouter.GET("/admin_menu/edit", admin.AdminGroupApp.AdminMenu.Edit)
		//菜单管理-更新菜单
		adminRouter.POST("/admin_menu/update", admin.AdminGroupApp.AdminMenu.Update)
		//菜单管理-删除菜单
		adminRouter.POST("/admin_menu/del", admin.AdminGroupApp.AdminMenu.Del)

		//操作日志
		adminRouter.GET("/admin_log/index", admin.AdminGroupApp.AdminLog.Index)

		//系统管理-开发管理-数据维护
		adminRouter.GET("/database/table", admin.AdminGroupApp.Database.Table)
		//系统管理-开发管理-数据维护-优化表
		adminRouter.POST("/database/optimize", admin.AdminGroupApp.Database.Optimize)
		//系统管理-开发管理-数据维护-修复表
		adminRouter.POST("/database/repair", admin.AdminGroupApp.Database.Repair)
		//系统管理-开发管理-数据维护-查看详情
		adminRouter.Match([]string{"GET", "POST"}, "/database/view", admin.AdminGroupApp.Database.View)

		//用户等级管理
		adminRouter.GET("/user_level/index", admin.AdminGroupApp.UserLevel.Index)
		//用户等级管理-添加界面
		adminRouter.GET("/user_level/add", admin.AdminGroupApp.UserLevel.Add)
		//用户等级管理-添加
		adminRouter.POST("/user_level/create", admin.AdminGroupApp.UserLevel.Create)
		//用户等级管理-修改界面
		adminRouter.GET("/user_level/edit", admin.AdminGroupApp.UserLevel.Edit)
		//用户等级管理-修改
		adminRouter.POST("/user_level/update", admin.AdminGroupApp.UserLevel.Update)
		//用户等级管理-启用
		adminRouter.POST("/user_level/enable", admin.AdminGroupApp.UserLevel.Enable)
		//用户等级管理-禁用
		adminRouter.POST("/user_level/disable", admin.AdminGroupApp.UserLevel.Disable)
		//用户等级管理-删除
		adminRouter.POST("/user_level/del", admin.AdminGroupApp.UserLevel.Del)
		//用户等级管理-导出
		adminRouter.GET("/user_level/export", admin.AdminGroupApp.UserLevel.Export)

		//用户管理
		adminRouter.GET("/user/index", admin.AdminGroupApp.User.Index)
		//用户管理-添加界面
		adminRouter.GET("/user/add", admin.AdminGroupApp.User.Add)
		//用户管理-添加
		adminRouter.POST("/user/create", admin.AdminGroupApp.User.Create)
		//用户管理-修改界面
		adminRouter.GET("/user/edit", admin.AdminGroupApp.User.Edit)
		//用户管理-修改
		adminRouter.POST("/user/update", admin.AdminGroupApp.User.Update)
		//用户管理-启用
		adminRouter.POST("/user/enable", admin.AdminGroupApp.User.Enable)
		//用户管理-禁用
		adminRouter.POST("/user/disable", admin.AdminGroupApp.User.Disable)
		//用户管理-删除
		adminRouter.POST("/user/del", admin.AdminGroupApp.User.Del)
		//用户管理-导出
		adminRouter.GET("/user/export", admin.AdminGroupApp.User.Export)

		//设置中心-后台设置
		adminRouter.GET("/setting/admin", admin.AdminGroupApp.Setting.Admin)
		//设置中心-更新设置
		adminRouter.POST("/setting/update", admin.AdminGroupApp.Setting.Update)
	}
}

func initSession() {

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: 7200, // 过期时间为 7200 秒
		Path:   "/",
	})
	router.Use(sessions.Sessions("mysession", store))
}

// 加载后台多模板
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/public/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(filepath.Dir(include))+"/"+filepath.Base(include), tmpFuncs, files...)
		fmt.Println("add include template:", filepath.Base(filepath.Dir(include))+"/"+filepath.Base(include))

	}
	r.AddFromFilesFuncs("database/view.html", tmpFuncs, "views/admin/includes/database/view.html", "views/admin/public/head_css.html", "views/admin/public/head_js.html")
	r.AddFromFilesFuncs("login.html", tmpFuncs, "views/admin/auth/login.html", "views/admin/auth/captcha.html")
	return r
}
