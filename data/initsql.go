package main

import (
	"BUAAJobHunting/model"
	"encoding/csv"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/yanyiwu/gojieba"
	"io"
	"log"
	"os"
	"strings"
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
		bar.Increment()
		row, err := r.Read()
		for i := 0; i < len(row); i++ {
			row[i] = strings.TrimSpace(row[i])
		}
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		sqlString := `INSERT INTO boss_data (job_name, job_area, salary, tag_list, hr_info, company_logo, 
            company_name, company_tag_list, company_url, job_need, job_desc, job_url, created_at, tokens) VALUES (
        	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10 , $11, $12, $13, $14)`
		if _, err := DB.Exec(sqlString, row[0], row[1], row[2], row[3], row[4], row[5],
			row[6], row[7], row[8], row[9], row[10], row[11], time.Now().Local(), ""); err != nil {
			panic(err)
		}
	}
	bar.Finish()
	fmt.Println("Finish loading " + fileName)
}

func InitBossToken() {
	sqlString := `SELECT * FROM boss_data`
	var bossJobs []model.BossJob
	if err := DB.Select(&bossJobs, sqlString); err != nil {
		panic(err)
	}
	x := gojieba.NewJieba()
	defer x.Free()
	var jobNameWords []string
	var jobNeedWords []string
	var jobDescWords []string
	var companyWords []string
	bar := pb.StartNew(1604)
	for _, bossJob := range bossJobs {
		bar.Increment()
		jobNameWords = x.CutForSearch(bossJob.JobName, true)
		jobNeedWords = x.CutForSearch(bossJob.JobNeed, true)
		jobDescWords = x.CutForSearch(bossJob.JobDesc, true)
		companyWords = x.CutForSearch(bossJob.CompanyName, true)
		sqlString = `UPDATE boss_data SET tokens = setweight(to_tsvector('simple', $1), 'A') || 
            setweight(to_tsvector('simple', $2), 'B') || setweight(to_tsvector('simple', $3), 'C') 
            || setweight(to_tsvector('simple', $4), 'D') WHERE id = $5`
		if _, err := DB.Exec(sqlString, strings.Join(jobNameWords, " "), strings.Join(jobNeedWords, " "),
			strings.Join(jobDescWords, " "), strings.Join(companyWords, " "), bossJob.ID); err != nil {
			panic(err)
		}
	}
	bar.Finish()
}

func Load58Data() {
	fileName := "./data/58_data.csv"
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
	bar := pb.StartNew(5630)
	for {
		bar.Increment()
		row, err := r.Read()
		for i := 0; i < len(row); i++ {
			row[i] = strings.TrimSpace(row[i])
		}
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		sqlString := `INSERT INTO "58_data" (job_name, job_area, salary, job_wel, company_name, 
            job_need, job_url, created_at, tokens) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		if _, err := DB.Exec(sqlString, row[0], row[1], row[2], row[3],
			row[4], row[5], row[6], time.Now().Local(), ""); err != nil {
			panic(err)
		}
	}
	bar.Finish()
	fmt.Println("Finish loading " + fileName)
}

func Init58Token() {
	sqlString := `SELECT * FROM "58_data"`
	var tc58Jobs []model.TC58Job
	if err := DB.Select(&tc58Jobs, sqlString); err != nil {
		panic(err)
	}
	x := gojieba.NewJieba()
	defer x.Free()
	var jobNameWords []string
	var jobNeedWords []string
	var jobWelWords []string
	bar := pb.StartNew(5630)
	for _, tc58Job := range tc58Jobs {
		bar.Increment()
		jobNameWords = x.CutForSearch(tc58Job.JobName, true)
		jobNeedWords = x.CutForSearch(tc58Job.JobNeed, true)
		jobWelWords = x.CutForSearch(tc58Job.JobWel, true)
		sqlString = `UPDATE "58_data" SET tokens = setweight(to_tsvector('simple', $1), 'A') || 
            setweight(to_tsvector('simple', $2), 'B') || setweight(to_tsvector('simple', $3), 'C') WHERE id = $4`
		if _, err := DB.Exec(sqlString, strings.Join(jobNameWords, " "), strings.Join(jobNeedWords, " "),
			strings.Join(jobWelWords, " "), tc58Job.ID); err != nil {
			panic(err)
		}
	}
	bar.Finish()
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
	InitBossToken()
	Load58Data()
	Init58Token()
}
