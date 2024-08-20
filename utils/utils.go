package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"gin-admin/global"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/bcrypt"
)

var TimeLayout = "2006-01-02 15:04:05"

func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

// CheckAndCreateDir checks if a directory exists and creates it if it doesn't
func CheckAndCreateDir(path string) error {
	// Check if the directory exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", path)
	} else if err != nil {
		// An error other than "directory does not exist" occurred
		return err
	} else if !info.IsDir() {
		// The path exists but it is not a directory
		return fmt.Errorf("%s already exists and is not a directory", path)
	} else {
		// The directory exists
		fmt.Println("Directory already exists:", path)
	}

	return nil
}

// TimeUntilTomorrowMidnight 计算从现在到明天0点的时间差
func TimeUntilTomorrowMidnight() time.Duration {
	now := time.Now()

	// 获取明天0点时间
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	// 计算时间差
	duration := tomorrow.Sub(now)

	return duration
}

// IsRunningUnderSystemd 判断当前进程是否运行在Systemd环境中
func IsRunningUnderSystemd() bool {
	if _, exists := os.LookupEnv("IS_RUNNING_SYSTEMD"); exists {
		return true
	}
	return false
}

func CreatePIDFile(pidFile string) error {
	pid := os.Getpid()
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	return err
}

// CaptchaResponse struct
type CaptchaResponse struct {
	CaptchaId  string
	CaptchaUrl string
}

// GetCaptcha 获取验证码
func GetCaptcha() *CaptchaResponse {
	captchaID := captcha.NewLen(4)
	return &CaptchaResponse{
		CaptchaId:  captchaID,
		CaptchaUrl: fmt.Sprintf("/admin/auth/captcha/%s", captchaID),
	}
}

// PasswordHash php的函数password_hash
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// PasswordVerify php的函数password_verify
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// KeyInMap 模仿php的array_key_exists,判断是否存在map中
func KeyInMap(key string, m map[string]int) bool {
	_, ok := m[key]
	if ok {
		return true
	}
	return false
}

// InArrayForInt 模仿php的in_array,判断是否存在int数组中
func InArrayForInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// InArrayForString 模仿php的in_array,判断是否存在string数组中
func InArrayForString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// IntArrToStringArr int数组转string数组
func IntArrToStringArr(arr []int) []string {
	var stringArr []string
	for _, v := range arr {
		stringArr = append(stringArr, strconv.Itoa(v))
	}
	return stringArr
}

// GetMd5String 对字符串进行MD5哈希
func GetMd5String(str string) string {
	t := md5.New()
	io.WriteString(t, str)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// GetSha1String 对字符串进行SHA1哈希
func GetSha1String(str string) string {
	t := sha1.New()
	io.WriteString(t, str)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// ParseName 字符串命名风格转换
func ParseName(name string, ptype int, ucfirst bool) string {
	if ptype > 0 {
		//解释正则表达式
		reg := regexp.MustCompile(`_([a-zA-Z])`)
		if reg == nil {
			global.LOG.Error("MustCompile err")
			return ""
		}
		//提取关键信息
		result := reg.FindAllStringSubmatch(name, -1)
		for _, v := range result {
			name = strings.ReplaceAll(name, v[0], strings.ToUpper(v[1]))
		}

		if ucfirst {
			return Ucfirst(name)
		}
		return Lcfirst(name)
	}
	//解释正则表达式
	reg := regexp.MustCompile(`[A-Z]`)
	if reg == nil {
		global.LOG.Error("MustCompile err")
		return ""
	}
	//提取关键信息
	result := reg.FindAllStringSubmatch(name, -1)

	for _, v := range result {
		name = strings.ReplaceAll(name, v[0], "_"+v[0])
	}
	return strings.ToLower(name)
}

// Ucfirst 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// Lcfirst 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
