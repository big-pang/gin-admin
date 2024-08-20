package main

import (
	"fmt"
	"gin-admin/initialize"

	"gin-admin/global"
	"gin-admin/global/logger"

	"gin-admin/router"
)

var ServiceName = "gin-admin"

func main() {
	initialize.Init()
	if err := router.Run(fmt.Sprintf("%s:%d", global.CONFIG.Base.Addr, global.CONFIG.Base.Port)); err != nil {
		logger.Sugar().Fatalln("服务器启动失败: %v", err)
	}
}
