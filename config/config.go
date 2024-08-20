package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Base struct {
		Name            string `yaml:"name"`
		Version         string `yaml:"version"`
		Addr            string `yaml:"addr"`
		Port            int    `yaml:"port"`
		Debug           bool   `yaml:"debug"`
		Mode            string `yaml:"mode"`
		LogAesKey       string `yaml:"log_aes_key"`
		ShortName       string `mapstructure:"short_name" json:"short_name" yaml:"short_name"`                   //后台简称
		Author          string `mapstructure:"author" json:"author" yaml:"author"`                               //后台作者
		Link            string `mapstructure:"link" json:"link" yaml:"link"`                                     //footer链接地址
		PasswordWarning string `mapstructure:"password_warning" json:"password_warning" yaml:"password_warning"` //默认密码警告
		ShowNotice      string `mapstructure:"show_notice" json:"show_notice" yaml:"show_notice"`                //是否显示提示信息
		NoticeContent   string `mapstructure:"notice_content" json:"notice_content" yaml:"notice_content"`       //提示信息内容
	} `yaml:"base"`

	Login  Login  `mapstructure:"login" json:"login" yaml:"login"`
	System System `mapstructure:"system" json:"system" yaml:"system"`

	Attachment Attachment `mapstructure:"attachment" json:"attachment" yaml:"attachment"`
	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
}

func InitConfig(configPath string) *Config {
	var config Config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件错误: %v", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("解析配置文件错误: %v", err)
	}
	return &config
}
