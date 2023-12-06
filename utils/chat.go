package utils

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"time"
)

func CheckChat() {
	sqlString := `SELECT * from last_chat WHERE time < $1`
	var lastChats []model.LastChat
	if err := global.Database.Select(&lastChats, sqlString, time.Now().In(time.FixedZone("CST", 8*3600-300))); err != nil {
		return
	}
	global.Dispatcher.Remind(lastChats)
}
