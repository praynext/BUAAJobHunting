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
	global.Router.GET("/all_user", GetAllUser)

	bossData := global.Router.Group("/boss_data")
	bossData.GET("/company", SearchBossDataByCompany)
	bossData.GET("/job", SearchBossDataByJob)
	bossData.GET("/favorite", UserGetFavoriteBossData)
	bossData.POST("/favorite", UserFavoriteBossData)
	bossData.DELETE("/favorite", UserCancelFavoriteBossData)

	tc58Data := global.Router.Group("/58_data")
	tc58Data.GET("/company", Search58DataByCompany)
	tc58Data.GET("/job", Search58DataByJob)
	tc58Data.GET("/favorite", UserGetFavorite58Data)
	tc58Data.POST("/favorite", UserFavorite58Data)
	tc58Data.DELETE("/favorite", UserCancelFavorite58Data)

	chat := global.Router.Group("/chat")
	chat.GET("/ws", ServeWebsocket)
	chat.GET("/history", GetChatHistory)

	reminder := global.Router.Group("/reminder")
	reminder.GET("/get", GetReminder)
	reminder.POST("/add", AddReminder)
	reminder.PUT("/update", UpdateReminder)
	reminder.DELETE("/delete", DeleteReminder)
}
