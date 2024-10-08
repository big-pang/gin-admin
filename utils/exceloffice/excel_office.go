package exceloffice

import (
	"gin-admin/global"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/gin-gonic/gin"
)

// ExportData 导出数据
func ExportData(head []string, body [][]string, name string, version string, title string, c *gin.Context) {
	if name == "" {
		name = time.Now().Format("2006-01-02-15-04-05")
	}

	if title == "" {
		title = "导出记录"
	}

	if version == "" {
		version = "2007"
	}

	charIndex := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	//处理超过26列
	a := "A"

	for _, v := range charIndex {
		charIndex = append(charIndex, a+v)
	}

	//创建excel
	f := excelize.NewFile()
	f.SetActiveSheet(0)

	//Excel表格头部
	for key, val := range head {
		f.SetCellValue("Sheet1", charIndex[key]+"1", val)
	}

	//Excel表格body部分
	for key, val := range body {
		row := key + 2
		col := 0
		for _, v := range val {
			f.SetCellValue("Sheet1", charIndex[col]+strconv.Itoa(row), v)
			col++
		}
	}

	//最后设置Sheet1自定义的title
	f.SetSheetName("Sheet1", title)

	//版本差异信息
	versionOpt := map[string]map[string]string{
		"2007": {
			"mime":       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"ext":        ".xlsx",
			"write_type": "Xlsx",
		},
		"2003": {
			"mime":       "application/vnd.ms-excel",
			"ext":        ".xls",
			"write_type": "Xls",
		},
		"pdf": {
			"mime":       "application/pdf",
			"ext":        ".pdf",
			"write_type": "PDF",
		},
		"ods": {
			"mime":       "application/vnd.oasis.opendocument.spreadsheet",
			"ext":        ".ods",
			"write_type": "OpenDocument",
		},
	}

	c.Header("Content-Type", versionOpt[version]["mime"])
	c.Header("Content-Disposition", "attachment;filename=\""+name+versionOpt[version]["ext"]+"\"")
	c.Header("Cache-Control", "max-age=0")

	if _, err := f.WriteTo(c.Writer); err != nil {
		global.LOG.Sugar().Warn("export data err:", err.Error())
	}
	return
}
