package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type JobResponse struct {
	JobName        string    `json:"job_name"`
	JobArea        string    `json:"job_area"`
	Salary         string    `json:"salary"`
	TagList        []string  `json:"tag_list"`
	HRInfo         []string  `json:"hr_info"`
	CompanyLogo    string    `json:"company_logo"`
	CompanyName    string    `json:"company_name"`
	CompanyTagList []string  `json:"company_tag_list"`
	CompanyURL     string    `json:"company_url"`
	JobNeed        []string  `json:"job_need"`
	JobDesc        string    `json:"job_desc"`
	JobURL         string    `json:"job_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type AllBossData struct {
	TotalCount int           `json:"total_count"`
	Jobs       []JobResponse `json:"jobs"`
}

// SearchBossDataByCompany godoc
// @Schemes http
// @Description 根据公司搜索Boss直聘数据
// @Tags Search
// @Param info query string true "用户搜索信息"
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllBossData "Boss直聘数据"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/company [get]
// @Security ApiKeyAuth
func SearchBossDataByCompany(c *gin.Context) {
	sqlString := `SELECT * FROM boss_data WHERE similarity($1, company_name) > 0.1 
        AND job_area LIKE $2 ORDER BY similarity($3, company_name) DESC, created_at DESC`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var bossJobs []model.BossJob
	if err := global.Database.Select(&bossJobs, sqlString, c.Query("info"),
		"%"+c.Query("area")+"%", c.Query("info")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []JobResponse
	for _, bossJob := range bossJobs {
		jobs = append(jobs, JobResponse{
			JobName:        bossJob.JobName,
			JobArea:        bossJob.JobArea,
			Salary:         bossJob.Salary,
			TagList:        strings.Split(bossJob.TagList, " "),
			HRInfo:         strings.Split(bossJob.HRInfo, " "),
			CompanyLogo:    bossJob.CompanyLogo,
			CompanyName:    bossJob.CompanyName,
			CompanyTagList: strings.Split(bossJob.CompanyTagList, " "),
			CompanyURL:     bossJob.CompanyURL,
			JobNeed:        strings.Split(bossJob.JobNeed, " "),
			JobDesc:        bossJob.JobDesc,
			JobURL:         bossJob.JobURL,
			CreatedAt:      bossJob.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, AllBossData{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

// SearchBossDataByJob godoc
// @Schemes http
// @Description 根据工作搜索Boss直聘数据
// @Tags Search
// @Param info query string true "用户搜索信息"
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllBossData "Boss直聘数据"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/job [get]
// @Security ApiKeyAuth
func SearchBossDataByJob(c *gin.Context) {
	sqlString := `SELECT * FROM boss_data WHERE tokens @@ to_tsquery('simple', $1) AND job_area 
    	LIKE $2 ORDER BY ts_rank(tokens, to_tsquery('simple', $3)) DESC, created_at DESC`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	queryWords := global.Parser.Cut(c.Query("info"), true)
	var bossJobs []model.BossJob
	if err := global.Database.Select(&bossJobs, sqlString, strings.Join(queryWords, " | "),
		"%"+c.Query("area")+"%", strings.Join(queryWords, " | ")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []JobResponse
	for _, bossJob := range bossJobs {
		jobs = append(jobs, JobResponse{
			JobName:        bossJob.JobName,
			JobArea:        bossJob.JobArea,
			Salary:         bossJob.Salary,
			TagList:        strings.Split(bossJob.TagList, " "),
			HRInfo:         strings.Split(bossJob.HRInfo, " "),
			CompanyLogo:    bossJob.CompanyLogo,
			CompanyName:    bossJob.CompanyName,
			CompanyTagList: strings.Split(bossJob.CompanyTagList, " "),
			CompanyURL:     bossJob.CompanyURL,
			JobNeed:        strings.Split(bossJob.JobNeed, " "),
			JobDesc:        bossJob.JobDesc,
			JobURL:         bossJob.JobURL,
			CreatedAt:      bossJob.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, AllBossData{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}
