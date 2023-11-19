package main

import (
	"BUAAJobHunting/api"
	"BUAAJobHunting/docs"
	"BUAAJobHunting/global"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yanyiwu/gojieba"
	"io"
	"log"
	"net/smtp"
	"os"
	"path"
)

func InitSql(Host string, Port int, User string, Password string, Database string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Host, Port, User, Password, Database)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	global.Database = db
}

func InitRedis(Host string, Port int, Password string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Host, Port),
		Password: Password,
		DB:       0,
	})
	global.Redis = rdb
}

func InitLog(Path string) {
	f, err := os.OpenFile(Path+"log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	ginLog, err := os.OpenFile(Path+"gin_log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	gin.DefaultWriter = io.MultiWriter(ginLog, os.Stdout)
}

func InitSMTP(Host string, User string, Password string) {
	global.SMTPAuth = smtp.PlainAuth("", User, Password, Host)
}

func InitParser(Path string) {
	jiebaPath := path.Join(Path, "jieba.dict.utf8")
	hmmPath := path.Join(Path, "hmm_model.utf8")
	userPath := path.Join(Path, "user.dict.utf8")
	idfPath := path.Join(Path, "idf.utf8")
	stopPath := path.Join(Path, "stop_words.utf8")
	global.Parser = gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	InitSql(viper.GetString("PostgresHost"), viper.GetInt("PostgresPort"),
		viper.GetString("PostgresUser"), viper.GetString("PostgresPassword"),
		viper.GetString("PostgresDatabase"))
	InitRedis(viper.GetString("RedisHost"), viper.GetInt("RedisPort"),
		viper.GetString("RedisPassword"))
	InitLog(viper.GetString("LogPath"))
	InitSMTP(viper.GetString("SMTPHost"), viper.GetString("SMTPUser"), viper.GetString("SMTPPassword"))
	InitParser(viper.GetString("DictPath"))
	docs.SwaggerInfo.BasePath = viper.GetString("DocsPath")
}

// @title BUAAJobHunting Backend API
// @version 0.0.1
// @license null
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Token
func main() {
	LoadConfig()
	global.Router = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, global.TokenHeader)
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Date")

	global.Router.Use(cors.New(corsConfig))
	global.Router.Use(global.Authenticate)
	global.Router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
	api.InitRoute()
	err := global.Router.Run("0.0.0.0:9000")
	if err != nil {
		return
	}
}
