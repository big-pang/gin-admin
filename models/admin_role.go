package models

// AdminRole struct
type AdminRole struct {
	Id          int    `gorm:"column(id);auto;size(11);description(表ID)" json:"id"`
	Name        string `gorm:"column(name);size(50);description(名称)" json:"name"`
	Description string `gorm:"column(description);size(100);description(简介)" json:"description"`
	Url         string `gorm:"column(url);size(1000);description(权限)" json:"url"`
	Status      int8   `gorm:"column(status);size(1);default(1);description(是否启用 0：否 1：是)" json:"status"`
}

// SearchField 定义模型的可搜索字段
func (*AdminRole) SearchField() []string {
	return []string{"name", "description"}
}

// NoDeletionId 禁止删除的数据id
func (*AdminRole) NoDeletionId() []int {
	return []int{1}
}

// WhereField 定义模型可作为条件的字段
func (*AdminRole) WhereField() []string {
	return []string{}
}

// TimeField 定义可做为时间范围查询的字段
func (*AdminRole) TimeField() []string {
	return []string{}
}

// TableName 自定义table 名称
func (*AdminRole) TableName() string {
	return "admin_role"
}
