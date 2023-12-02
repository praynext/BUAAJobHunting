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
			Time: msg.Time.Format("2006/01/02 15:04:05"),
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

type AllMessageData struct {
	TotalCount int              `json:"total_count"`
	Messages   []global.Message `json:"messages"`
}

// GetChatHistory godoc
// @Schemes http
// @Description 获取用户聊天记录
// @Tags Chat
// @Param user_id query int true "用户id"
// @Param offset query int false "偏移量"
// @Param limit query int false "限制条数"
// @Success 200 {object} AllMessageData "用户聊天记录"
// @Failure default {string} string "服务器错误"
// @Router /chat/history [get]
// @Security ApiKeyAuth
func GetChatHistory(c *gin.Context) {
	sqlString := `SELECT * FROM "message" WHERE ("from"=$1 AND "to"=$2) OR ("to"=$3 AND "from"=$4) ORDER BY time DESC`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var messages []model.Message
	if err := global.Database.Select(&messages, sqlString, c.Query("user_id"),
		c.GetInt("UserId"), c.Query("user_id"), c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var messageResponse []global.Message
	for _, msg := range messages {
		messageResponse = append(messageResponse, global.Message{
			From: msg.From,
			To:   msg.To,
			Msg:  msg.Msg,
			Time: msg.Time.Format("2006/01/02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, AllMessageData{
		TotalCount: len(messageResponse),
		Messages:   messageResponse,
	})
}
