package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
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

	var messages []model.Message
	sqlString := `SELECT * FROM "message" WHERE "to"=$1 AND has_sent=false`
	if err := global.Database.Select(&messages, sqlString, client.UserId); err != nil {
		log.Fatalf("Query messages failed: %v", err)
		return
	}
	for _, msg := range messages {
		payload := global.Message{
			From: msg.From,
			To:   msg.To,
			Msg:  msg.Msg,
			Time: msg.Time,
		}
		bytePayload, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Marshal messages failed: %v", err)
			return
		}
		client.Send <- bytePayload

		// Save msg into database, setting has_sent true
		_, err = global.Database.Exec(`UPDATE "message" SET has_sent=true WHERE "from"=$1 AND "to"=$2 AND message=$3`, msg.From, client.UserId, msg.Msg)
		if err != nil {
			log.Fatalf("Save message into database failed: %v", err)
			return
		}
	}

	go client.WritePump()
	go client.ReadPump()
}
