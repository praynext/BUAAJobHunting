package main

import (
	"encoding/csv"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"time"
)

var DB *sqlx.DB

func InitSql(Host string, Port int, User string, Password string, Database string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Host, Port, User, Password, Database)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	DB = db
}

func LoadBossData() {
	fileName := "./data/boss_data.csv"
	fs, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(fs *os.File) {
		err := fs.Close()
		if err != nil {
			log.Fatalf("can not close, err is %+v", err)
		}
	}(fs)
	r := csv.NewReader(fs)
	bar := pb.StartNew(1604)
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		bar.Increment()
		sqlString := `INSERT INTO boss_data (job_name, job_area, salary, tag_list, hr_info, company_logo, 
            company_name, company_tag_list, company_url, job_need, job_desc, job_url,created_at) VALUES (
        	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10 , $11, $12, $13)`
		if _, err := DB.Exec(sqlString, row[0], row[1], row[2], row[3], row[4], row[5],
			row[6], row[7], row[8], row[9], row[10], row[11], time.Now().Local()); err != nil {
			panic(err)
		}
	}
	bar.Finish()
	fmt.Println("Finish loading " + fileName)
}

func main() {
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
	LoadBossData()
}
