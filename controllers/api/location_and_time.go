package api

import (
	"gin-admin/utils/ipsearch"

	"github.com/gin-gonic/gin"
)

func LocationAndTime(c *gin.Context) {

	ip := c.RemoteIP()
	// ip = "210.51.200.123"
	location := ipsearch.IpSearch.GetLocation(ip)
	c.PureJSON(200, location)
}
