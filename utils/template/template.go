// Package template 自定义模板函数
package template

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"time"
)

// UnixTimeForFormat 时间轴转时间字符串
func UnixTimeForFormat(timeUnix int) string {
	//转化所需模板
	timeLayout := "2006-01-02 15:04:05"
	return time.Unix(int64(timeUnix), 0).Format(timeLayout)
}

// TimeForFormat 时间转时间字符串
func TimeForFormat(t time.Time) string {
	//转化所需模板
	timeLayout := "2006-01-02 15:04:05"
	return t.Format(timeLayout)
}

// FormatSize 格式化文件大小单位
func FormatSize(size string, delimiter string) string {
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return ""
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; sizeInt >= 1024 && i < 5; i++ {
		sizeInt /= 1024
	}
	return strconv.FormatFloat(math.Round(float64(sizeInt)), 'f', -1, 64) + delimiter + units[i]
}

func Str2Html(s string) template.HTML {
	return template.HTML(s)
}

func AssetsCSS(path string) template.HTML {
	return template.HTML("<link rel=\"stylesheet\" href=\"/assets/" + path + "\">")
}

func AssetsJS(path string) template.HTML {
	return template.HTML("<script src=\"/assets/" + path + "\"></script>")
}

func Compare(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func MapGet(m interface{}, key interface{}) interface{} {
	switch v := m.(type) {
	case map[string]interface{}:
		return v[key.(string)]
	case map[int]string:
		return v[key.(int)]
	default:
		return nil
	}
}
