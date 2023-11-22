package api

import (
	"BUAAJobHunting/global"
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

	// send all has_sent=false messages to client
	exec, err := global.Database.Query(`SELECT "from", message FROM "message" WHERE "to"=$1 AND has_sent=false`, client.UserId)
	if err != nil {
		log.Fatalf("Query messages failed: %v", err)
		return
	}
	for exec.Next() {
		var from int
		var msg string
		if err := exec.Scan(&from, &msg); err != nil {
			log.Fatalf("Scan messages failed: %v", err)
			return
		}
		payload := global.Message{
			From: from,
			To:   client.UserId,
			Msg:  msg,
		}
		byte_payload, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Marshal messages failed: %v", err)
			return
		}
		client.Send <- byte_payload

		// Save msg into database, setting has_sent true
		_, err = global.Database.Exec(`UPDATE "message" SET has_sent=true WHERE "from"=$1 AND "to"=$2 AND message=$3`, from, client.UserId, msg)
		if err != nil {
			log.Fatalf("Save message into database failed: %v", err)
			return
		}
	}

	go client.WritePump()
	go client.ReadPump()
}
