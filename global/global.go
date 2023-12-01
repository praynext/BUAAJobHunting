package global

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/yanyiwu/gojieba"
	"net/smtp"
)

var Database *sqlx.DB
var Redis *redis.Client
var Router *gin.Engine
var SMTPAuth smtp.Auth
var Parser *gojieba.Jieba
var Dispatcher *Hub
var UpGrader *websocket.Upgrader
var Cron *cron.Cron
