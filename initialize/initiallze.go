package initialize

import (
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"gin-admin/global/logger"
	"gin-admin/router"
	"gin-admin/utils/ipsearch"
)

// 显示的初始化
func Init() {
	fmt.Println("初始化程序......")
	// 首先初始化配置
	global.CONFIG = config.InitConfig("config.yaml")
	// 初始化日志
	logger.LogInit("gin-admin")
	// 初始化数据库
	global.DB = Gorm()
	// 初始化路由Init
	router.InitRouter()
	ipsearch.Init()
}

func init() {
	Init()
}
