package models

import (
	"time"

	"gorm.io/gorm"
)

// User struct
type User struct {
	Id          int            `gorm:"column(id);auto;size(11);description(表ID)" json:"id"`
	Avatar      string         `gorm:"column(avatar);size(255);description(头像)" json:"avatar"`
	Username    string         `gorm:"column(username);size(30);description(用户名)" json:"username"`
	Nickname    string         `gorm:"column(nickname);size(30);description(昵称)" json:"nickname"`
	Mobile      string         `gorm:"column(mobile);size(11);description(手机号)" json:"mobile"`
	UserLevelId int            `gorm:"column(user_level_id);size(11);default(1);description(用户等级)" json:"user_level_id"`
	Password    string         `gorm:"column(password);size(255);description(密码)" json:"password"`
	Status      int8           `gorm:"column(status);size(1);default(1);description(是否启用 0：否 1：是)" json:"status"`
	Description string         `gorm:"column(description);type(text);description(描述)" json:"description"`
	CreatedAt   time.Time      `gorm:"column(created_at);type(timestamp);default(CURRENT_TIMESTAMP);description(创建时间)" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column(updated_at);type(timestamp);default(CURRENT_TIMESTAMP);description(更新时间)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `orm:"column(deleted_at);type(timestamp);default(NULL);description(删除时间)" json:"deleted_at"`
}

// TableName 自定义table 名称
func (*User) TableName() string {
	return "user"
}

// SearchField 定义模型的可搜索字段
func (*User) SearchField() []string {
	return []string{"username", "mobile", "nickname"}
}

// NoDeletionId 禁止删除的数据id
func (*User) NoDeletionId() []int {
	return []int{}
}

// WhereField 定义模型可作为条件的字段
func (*User) WhereField() []string {
	return []string{}
}

// TimeField 定义可做为时间范围查询的字段
func (*User) TimeField() []string {
	return []string{}
}
