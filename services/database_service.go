package services

import (
	"fmt"
	"gin-admin/global"
)

// DbVersion struct
type DbVersion struct {
	DbVersion string
}

// DatabaseService struct
type DatabaseService struct {
	baseService
}

// GetMysqlVersion 获取mysql的版本
func (*DatabaseService) GetMysqlVersion() string {
	var dbVersion DbVersion
	err := global.DB.Raw("select VERSION() as db_version").First(&dbVersion).Error
	if err != nil {
		return "not found."
	}
	return dbVersion.DbVersion
}

// GetTableStatus 获取所有数据表的状态
func (ds *DatabaseService) GetTableStatus() ([]map[string]string, int) {
	var maps []map[string]interface{}
	var resultMaps []map[string]string

	result := global.DB.Raw("SHOW TABLE STATUS").Find(&maps)

	if result.RowsAffected > 0 && result.Error == nil {
		for _, item := range maps {
			resultMaps = append(resultMaps, map[string]string{
				"name":        ds.nil2String(item["Name"]),
				"comment":     ds.nil2String(item["Comment"]),
				"engine":      ds.nil2String(item["Engine"]),
				"collation":   ds.nil2String(item["Collation"]),
				"data_length": ds.nil2String(item["Data_length"]),
				"created_at":  ds.nil2String(item["Create_time"]),
				"updated_at":  ds.nil2String(item["Update_time"]),
			})
		}
	}

	return resultMaps, int(result.RowsAffected)
}

// OptimizeTable 优化数据表
func (*DatabaseService) OptimizeTable(tableName string) bool {
	err := global.DB.Exec("OPTIMIZE TABLE `" + tableName + "`").Error
	if err == nil {
		return true
	}
	return false
}

// RepairTable 修复数据表
func (*DatabaseService) RepairTable(tableName string) bool {
	err := global.DB.Exec("REPAIR TABLE `" + tableName + "`").Error
	if err == nil {
		return true
	}
	return false
}

// GetFullColumnsFromTable 获取数据表的所有字段
func (ds *DatabaseService) GetFullColumnsFromTable(tableName string) []map[string]string {
	var maps []map[string]interface{}
	var resultMaps []map[string]string
	result := global.DB.Raw("SHOW FULL COLUMNS FROM `" + tableName + "`").Find(&maps)

	if result.RowsAffected > 0 && result.Error == nil {
		for _, item := range maps {
			resultMaps = append(resultMaps, map[string]string{
				"name":       ds.nil2String(item["Field"]),
				"type":       ds.nil2String(item["Type"]),
				"collation":  ds.nil2String(item["Collation"]),
				"null":       ds.nil2String(item["Null"]),
				"key":        ds.nil2String(item["Key"]),
				"default":    ds.nil2String(item["Default"]),
				"extra":      ds.nil2String(item["Extra"]),
				"privileges": ds.nil2String(item["Privileges"]),
				"comment":    ds.nil2String(item["Comment"]),
			})
		}
	}

	return resultMaps
}

// nil2String interface 转换 为string
func (*DatabaseService) nil2String(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v) // 将其他类型转换为字符串
	}
}
