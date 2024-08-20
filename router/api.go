package router

import (
	apiControlls "gin-admin/controllers/api"
)

func apiRouter() {
	// 设置中间件
	// r.Use(middleware.AuthMiddle())
	api := router.Group("/api")
	{
		api.Any("/location_and_time", apiControlls.LocationAndTime)
	}
}
