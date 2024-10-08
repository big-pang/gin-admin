package models

// AdminLogData struct
type AdminLogData struct {
	Id         int    `gorm:"column(id);auto;size(11);description(表ID)" json:"id"`
	AdminLogId int    `gorm:"column(admin_log_id);size(11);description(日志ID)" json:"admin_log_id"`
	Data       string `gorm:"column(data);type(text);description(日志内容)" json:"data"`
}

// TableName 自定义table 名称
func (*AdminLogData) TableName() string {
	return "admin_log_data"
}

// SearchField 定义模型的可搜索字段
func (*AdminLogData) SearchField() []string {
	return []string{}
}

// WhereField 定义模型可作为条件的字段
func (*AdminLogData) WhereField() []string {
	return []string{}
}

// TimeField 定义可做为时间范围查询的字段
func (*AdminLogData) TimeField() []string {
	return []string{}
}

// NoDeletionId 禁止删除的数据id
func (*AdminLogData) NoDeletionId() []int {
	return []int{}
}
