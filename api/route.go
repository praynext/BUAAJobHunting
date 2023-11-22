package api

import (
	"BUAAJobHunting/global"
)

func InitRoute() {
	global.Router.GET("/ping", Ping)
	global.Router.POST("/login", Login)
	global.Router.POST("/logout", Logout)
	global.Router.POST("/register", Register)
	global.Router.POST("/change_password", ChangePassword)
	global.Router.POST("/reset_password", ResetPassword)
	global.Router.POST("/send_email", SendEmail)

	bossData := global.Router.Group("/boss_data")
	bossData.GET("/company", SearchBossDataByCompany)
	bossData.GET("/job", SearchBossDataByJob)

	tc58Data := global.Router.Group("/58_data")
	tc58Data.GET("/company", Search58DataByCompany)
	tc58Data.GET("/job", Search58DataByJob)

	chat := global.Router.Group("/chat")
	chat.GET("/ws", ServeWebsocket)
}
