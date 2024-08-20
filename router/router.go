package router

import (
	"gin-admin/global"
	"gin-admin/global/logger"
	templateFunc "gin-admin/utils/template"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var router *gin.Engine
var namedRoutes *NamedRoutes
var tmpFuncs template.FuncMap

func InitRouter() {

	if global.CONFIG.Base.Mode == "release" {
		router = gin.New()
		router.Use(GinLogger(logger.Logger), GinRecovery(logger.Logger, true))
	} else {
		router = gin.Default()
	}
	tmpFuncs = template.FuncMap{
		"UnixTimeForFormat": templateFunc.UnixTimeForFormat,
		"TimeForFormat":     templateFunc.TimeForFormat,
		"FormatSize":        templateFunc.FormatSize,
		"Urlfor":            namedRoutes.URLFor,
		"str2html":          templateFunc.Str2Html,
		"assets_css":        templateFunc.AssetsCSS,
		"assets_js":         templateFunc.AssetsJS,
		"compare":           templateFunc.Compare,
		"map_get":           templateFunc.MapGet,
	}
	router.SetFuncMap(tmpFuncs)
	apiRouter()
	adminRouter()
}

func Run(addr ...string) error {
	return router.Run(addr...)
}

func GinLogger(lgoger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		)
	}
}

func GinRecovery(lgoger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)

				// 返回 500 状态码
				if !c.IsAborted() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "Internal Server Error",
					})
				}
			}

		}()
		c.Next()
	}
}
