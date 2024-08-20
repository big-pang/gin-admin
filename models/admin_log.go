package models

import (
	"time"
)

// AdminLog struct

type AdminLog struct {
	Id          int       `gorm:"primaryKey;autoIncrement;column:id"`
	AdminUserId int       `gorm:"not null;column:admin_user_id" comment:"用户"`
	Name        string    `gorm:"type:varchar(30);not null;default:'';column:name" comment:"操作"`
	Url         string    `gorm:"type:varchar(100);not null;default:'';column:url" comment:"URL"`
	LogMethod   string    `gorm:"type:varchar(8);not null;default:'不记录';column:log_method" comment:"记录日志方法"`
	LogIp       string    `gorm:"type:varchar(20);not null;default:'';column:log_ip" comment:"操作IP"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;column:created_at" comment:"操作时间"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" comment:"更新时间"`
}

// SearchField 定义模型的可搜索字段
func (*AdminLog) SearchField() []string {
	return []string{"name", "url", "log_ip"}
}

// WhereField 定义模型可作为条件的字段
func (*AdminLog) WhereField() []string {
	return []string{"admin_user_id"}
}

// TimeField 定义可做为时间范围查询的字段
func (*AdminLog) TimeField() []string {
	return []string{"created_at"}
}

// NoDeletionId 禁止删除的数据id
func (*AdminLog) NoDeletionId() []int {
	return []int{}
}

// TableName 自定义table 名称
func (*AdminLog) TableName() string {
	return "admin_log"
}
