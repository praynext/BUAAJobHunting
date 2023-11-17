package global

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"net/smtp"
)

var Database *sqlx.DB
var Redis *redis.Client
var Router *gin.Engine
var SMTPAuth smtp.Auth
