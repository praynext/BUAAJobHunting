package api

import "BUAAJobHunting/global"

func InitRoute() {
	global.Router.GET("/ping", Ping)
	global.Router.POST("/login", Login)
	global.Router.POST("/logout", Logout)
	global.Router.POST("/register", Register)
	global.Router.POST("/change_password", ChangePassword)
	global.Router.POST("/reset_password", ResetPassword)
	global.Router.POST("/send_email", SendEmail)
}
