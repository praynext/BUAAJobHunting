package api

import (
	"BUAAJobHunting/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServeWebsocket(c *gin.Context) {
	conn, err := global.UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "建立websocket失败")
		return
	}
	client := &global.Client{
		UserId: c.GetInt("UserId"),
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	global.Dispatcher.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
