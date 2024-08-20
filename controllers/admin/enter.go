package admin

type AdminGroup struct {
	Auth
	Index
	AdminUser
	AdminRole
	AdminMenu
	AdminLog
	Database
	User
	UserLevel
	Setting
}

var AdminGroupApp = new(AdminGroup)
