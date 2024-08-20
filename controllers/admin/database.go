package admin

import (
	"gin-admin/global/response"
	"gin-admin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Database struct
type Database struct {
	base
}

// Table 显示数据表
func (dc *Database) Table(c *gin.Context) {
	dc.getBaseData(c)
	var databaseService services.DatabaseService
	data, affectRows := databaseService.GetTableStatus()

	dc.data["data"] = data
	dc.data["total"] = affectRows

	c.HTML(http.StatusOK, "database/table.html", dc.data)
}

// Optimize 优化表
func (dc *Database) Optimize(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		response.ErrorWithMessage("请指定要优化的表", c)
		return
	}
	var databaseService services.DatabaseService
	ok := databaseService.OptimizeTable(name)
	if ok {
		response.SuccessWithMessage("数据表"+name+"优化成功", c)
	} else {
		response.ErrorWithMessage("数据表"+name+"优化失败", c)
	}
}

// Repair 修复数据表
func (dc *Database) Repair(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		response.ErrorWithMessage("请指定要修复的表", c)
		return
	}
	var databaseService services.DatabaseService
	ok := databaseService.OptimizeTable(name)
	if ok {
		response.SuccessWithMessage("数据表"+name+"修复成功", c)
	} else {
		response.ErrorWithMessage("数据表"+name+"修复失败", c)
	}
}

// View 查看数据表
func (dc *Database) View(c *gin.Context) {
	dc.getBaseData(c)
	name := c.Query("name")

	if name == "" {
		response.ErrorWithMessage("请指定要查看的表", c)
	}

	var databaseService services.DatabaseService
	data := databaseService.GetFullColumnsFromTable(name)

	dc.data["data"] = data

	c.HTML(http.StatusOK, "database/view.html", dc.data)
}
