package response

import (
	"gin-admin/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR   = 0
	SUCCESS = 1
)

// Response 响应参数结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Url  string      `json:"url"`
	Wait int         `json:"wait"`
}

// Result 返回结果辅助函数
func Result(code int, msg string, data interface{}, url string, wait int, header map[string]string, c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		result := Response{
			Code: code,
			Msg:  msg,
			Data: data,
			Url:  url,
			Wait: wait,
		}

		if len(header) > 0 {
			for k, v := range header {
				c.Header(k, v)
			}
		}

		c.JSON(http.StatusOK, result)

		// 中断当前的请求
		c.Abort()
		return
	}

	if url == "" {
		url = c.Request.Referer()
		if url == "" {
			url = "/admin/index/index"
		}
	}

	c.Redirect(http.StatusFound, url)
}

// Success 成功、普通返回
func Success(c *gin.Context) {
	Result(SUCCESS, "操作成功", "", global.URL_BACK, 0, map[string]string{}, c)
}

// SuccessWithMessage 成功、返回自定义信息
func SuccessWithMessage(msg string, c *gin.Context) {
	Result(SUCCESS, msg, "", global.URL_BACK, 0, map[string]string{}, c)
}

// SuccessWithMessageAndUrl 成功、返回自定义信息和url
func SuccessWithMessageAndUrl(msg string, url string, c *gin.Context) {
	Result(SUCCESS, msg, "", url, 0, map[string]string{}, c)
}

// SuccessWithDetailed 成功、返回所有自定义信息
func SuccessWithDetailed(msg string, url string, data interface{}, wait int, header map[string]string, c *gin.Context) {
	Result(SUCCESS, msg, data, url, wait, header, c)
}

// Error 失败、普通返回
func Error(c *gin.Context) {
	Result(ERROR, "操作失败", "", global.URL_CURRENT, 0, map[string]string{}, c)
}

// ErrorWithMessage 失败、返回自定义信息
func ErrorWithMessage(msg string, c *gin.Context) {
	Result(ERROR, msg, "", global.URL_CURRENT, 0, map[string]string{}, c)
}

// ErrorWithMessageAndUrl 失败、返回自定义信息和url
func ErrorWithMessageAndUrl(msg string, url string, c *gin.Context) {
	Result(ERROR, msg, "", url, 0, map[string]string{}, c)
}

// ErrorWithDetailed 失败、返回所有自定义信息
func ErrorWithDetailed(msg string, url string, data interface{}, wait int, header map[string]string, c *gin.Context) {
	Result(ERROR, msg, data, url, wait, header, c)
}
