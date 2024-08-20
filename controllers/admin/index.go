package admin

import (
	"bufio"
	"encoding/base64"

	"gin-admin/global"
	"gin-admin/services"
	"gin-admin/utils"

	"gin-admin/utils/ipsearch"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Index struct{}

type PackageLib struct {
	Name    string
	Version string
}

func (ic *Index) Index(c *gin.Context) {
	data := gin.H{}
	loginUser, baseData, _ := getAdminData(c)
	data["login_user"] = loginUser

	//默认密码修改检测
	data["password_danger"] = 0

	//是否首页显示提示信息
	data["show_notice"] = global.CONFIG.Base.ShowNotice
	//提示内容
	data["notice_content"] = global.CONFIG.Base.NoticeContent

	//默认密码修改检测
	loginUserPassword, _ := base64.StdEncoding.DecodeString(loginUser.Password)
	if global.CONFIG.Base.PasswordWarning == "1" && utils.PasswordVerify("123456", string(loginUserPassword)) {
		data["password_danger"] = 1
	}

	//后台用户数量
	var adminUserService services.AdminUserService
	data["admin_user_count"] = adminUserService.GetCount()
	//后台角色数量
	var adminRoleService services.AdminRoleService
	data["admin_role_count"] = adminRoleService.GetCount()
	//后台菜单数量
	var adminMenuService services.AdminMenuService
	data["admin_menu_count"] = adminMenuService.GetCount()
	//后台日志数量
	var adminLogService services.AdminLogService
	data["admin_log_count"] = adminLogService.GetCount()
	//系统信息
	data["system_info"] = ic.getSystemInfo(c)
	// data["content"] = "views/admin/index/index.html"

	for k, v := range baseData {
		data[k] = v
	}

	c.HTML(200, "index/index.html", data)

}

// getSystemInfo 获取系统信息
func (ic *Index) getSystemInfo(c *gin.Context) map[string]interface{} {
	systemInfo := make(map[string]interface{})
	//服务器系统
	systemInfo["server_os"] = runtime.GOOS
	//Go版本
	systemInfo["go_version"] = runtime.Version()
	//文件上传默认内存缓存大小
	// systemInfo["upload_file_max_memory"] = int(beego.BConfig.MaxMemory / 1024 / 1024)
	//gin版本
	systemInfo["gin_version"] = gin.Version
	//当前后台版本
	systemInfo["admin_version"] = global.CONFIG.Base.Version
	//mysql版本
	var databaseService services.DatabaseService
	systemInfo["db_version"] = databaseService.GetMysqlVersion()
	//go时区
	systemInfo["timezone"] = time.UTC
	//当前时间
	systemInfo["date_time"] = time.Now().Format("2006-01-02 15:04:05")
	//用户IP
	systemInfo["user_ip"] = c.ClientIP()

	city := ipsearch.IpSearch.GetLocation(c.ClientIP()).City
	systemInfo["user_city"] = city

	userAgent := c.Request.Header.Get("user-agent")

	userOs := "Other"
	if strings.Contains(userAgent, "win") {
		userOs = "Windows"
	} else if strings.Contains(userAgent, "mac") {
		userOs = "MAC"
	} else if strings.Contains(userAgent, "linux") {
		userOs = "Linux"
	} else if strings.Contains(userAgent, "unix") {
		userOs = "Unix"
	} else if strings.Contains(userAgent, "bsd") {
		userOs = "BSD"
	} else if strings.Contains(userAgent, "iPad") || strings.Contains(userAgent, "iPhone") {
		userOs = "IOS"
	} else if strings.Contains(userAgent, "android") {
		userOs = "Android"
	}

	userBrowser := "Other"
	if strings.Contains(userAgent, "MSIE") {
		userBrowser = "MSIE"
	} else if strings.Contains(userAgent, "Firefox") {
		userBrowser = "Firefox"
	} else if strings.Contains(userAgent, "Chrome") {
		userBrowser = "Chrome"
	} else if strings.Contains(userAgent, "Safari") {
		userBrowser = "Safari"
	} else if strings.Contains(userAgent, "Opera") {
		userBrowser = "Opera"
	}

	//用户系统
	systemInfo["user_os"] = userOs
	//用户浏览器
	systemInfo["user_browser"] = userBrowser

	//读取go.mod文件
	var requireList []*PackageLib
	srcFile, err := os.Open("go.mod")
	if err != nil {
		global.LOG.Sugar().Error(err)
	} else {
		defer srcFile.Close()
		reader := bufio.NewReader(srcFile)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}
			if string(line) != "" {
				strArr := strings.Split(strings.TrimSpace(string(line)), " ")
				lenStrArr := len(strArr)
				//常规require方式
				if strArr[0] == "require" && lenStrArr >= 3 {
					packageLib := PackageLib{
						Name:    strArr[1],
						Version: strArr[2],
					}
					requireList = append(requireList, &packageLib)
				} else {
					//require多个时候
					if lenStrArr >= 2 && strings.Contains(strArr[0], "/") {
						packageLib := PackageLib{
							Name:    strArr[0],
							Version: strArr[1],
						}
						requireList = append(requireList, &packageLib)
					}
				}
			}
		}
	}

	systemInfo["require_list"] = requireList

	return systemInfo
}
