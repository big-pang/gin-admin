package models

import "encoding/gob"

// Params stores the Params
type Params map[string]interface{}

type ParamsId struct {
	Id int64 `json:"id"`
}

// ParamsList stores paramslist
type ParamsList []interface{}

func init() {
	gob.Register(AdminUser{})
	gob.Register(AdminLog{})
	gob.Register(AdminLogData{})
	gob.Register(AdminMenu{})
	gob.Register(AdminRole{})
	gob.Register(Attachment{})
	gob.Register(User{})
	gob.Register(UserLevel{})
	gob.Register(Setting{})
}

func menuArrToMap(menuList map[int]AdminMenu) map[int]Params {
	returnMaps := make(map[int]Params)

	for k, menu := range menuList {
		returnMaps[k] = Params{
			"Id":       menu.Id,
			"ParentId": menu.ParentId,
			"Name":     menu.Name,
			"Url":      menu.Url,
			"Icon":     menu.Icon,
			"SortId":   menu.SortId,
		}

	}
	return returnMaps
}
