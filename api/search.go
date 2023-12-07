package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type BossJobResponse struct {
	JobId          int       `json:"job_id"`
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
	IsFull         bool      `json:"is_full"`
	IsFavor        bool      `json:"is_favor"`
}

type AllBossData struct {
	TotalCount int               `json:"total_count"`
	Jobs       []BossJobResponse `json:"jobs"`
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
	var jobs []BossJobResponse
	for _, bossJob := range bossJobs {
		sqlString = `SELECT count(*) FROM user_favorite_boss_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), bossJob.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, BossJobResponse{
			JobId:          bossJob.ID,
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
			IsFull:         bossJob.IsFull,
			IsFavor:        isFavor == 1,
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
	var jobs []BossJobResponse
	for _, bossJob := range bossJobs {
		sqlString = `SELECT count(*) FROM user_favorite_boss_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), bossJob.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, BossJobResponse{
			JobId:          bossJob.ID,
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
			IsFull:         bossJob.IsFull,
			IsFavor:        isFavor == 1,
		})
	}
	c.JSON(http.StatusOK, AllBossData{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

// SearchBossDataByRandom godoc
// @Schemes http
// @Description 随机搜索Boss直聘数据
// @Tags Search
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllBossData "Boss直聘数据"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/random [get]
// @Security ApiKeyAuth
func SearchBossDataByRandom(c *gin.Context) {
	sqlString := `SELECT * FROM boss_data WHERE job_area LIKE $1 ORDER BY random()`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var bossJobs []model.BossJob
	if err := global.Database.Select(&bossJobs, sqlString, "%"+c.Query("area")+"%"); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []BossJobResponse
	for _, bossJob := range bossJobs {
		sqlString = `SELECT count(*) FROM user_favorite_boss_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), bossJob.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, BossJobResponse{
			JobId:          bossJob.ID,
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
			IsFull:         bossJob.IsFull,
			IsFavor:        isFavor == 1,
		})
	}
	c.JSON(http.StatusOK, AllBossData{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

type TC58JobResponse struct {
	JobId       int       `json:"job_id"`
	JobName     string    `json:"job_name"`
	JobArea     string    `json:"job_area"`
	Salary      string    `json:"salary"`
	JobWel      []string  `json:"job_wel"`
	CompanyName string    `json:"company_name"`
	JobNeed     []string  `json:"job_need"`
	JobURL      string    `json:"job_url"`
	CreatedAt   time.Time `json:"created_at"`
	IsFull      bool      `json:"is_full"`
	IsFavor     bool      `json:"is_favor"`
}

type All58Data struct {
	TotalCount int               `json:"total_count"`
	Jobs       []TC58JobResponse `json:"jobs"`
}

// Search58DataByCompany godoc
// @Schemes http
// @Description 根据公司搜索58同城数据
// @Tags Search
// @Param info query string true "用户搜索信息"
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} All58Data "58同城数据"
// @Failure default {string} string "服务器错误"
// @Router /58_data/company [get]
// @Security ApiKeyAuth
func Search58DataByCompany(c *gin.Context) {
	sqlString := `SELECT * FROM "58_data" WHERE similarity($1, company_name) > 0.1 
        AND job_area LIKE $2 ORDER BY similarity($3, company_name) DESC, created_at DESC`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var tc58Jobs []model.TC58Job
	if err := global.Database.Select(&tc58Jobs, sqlString, c.Query("info"),
		"%"+c.Query("area")+"%", c.Query("info")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []TC58JobResponse
	for _, tc58Job := range tc58Jobs {
		sqlString = `SELECT count(*) FROM user_favorite_58_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), tc58Job.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, TC58JobResponse{
			JobId:       tc58Job.ID,
			JobName:     tc58Job.JobName,
			JobArea:     tc58Job.JobArea,
			Salary:      tc58Job.Salary,
			JobWel:      strings.Split(tc58Job.JobWel, " "),
			CompanyName: tc58Job.CompanyName,
			JobNeed:     strings.Split(tc58Job.JobNeed, " "),
			JobURL:      tc58Job.JobURL,
			CreatedAt:   tc58Job.CreatedAt,
			IsFull:      tc58Job.IsFull,
			IsFavor:     isFavor == 1,
		})
	}
	c.JSON(http.StatusOK, All58Data{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

// Search58DataByJob godoc
// @Schemes http
// @Description 根据工作搜索58同城数据
// @Tags Search
// @Param info query string true "用户搜索信息"
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} All58Data "58同城数据"
// @Failure default {string} string "服务器错误"
// @Router /58_data/job [get]
// @Security ApiKeyAuth
func Search58DataByJob(c *gin.Context) {
	sqlString := `SELECT * FROM "58_data" WHERE tokens @@ to_tsquery('simple', $1) AND job_area 
    	LIKE $2 ORDER BY ts_rank(tokens, to_tsquery('simple', $3)) DESC, created_at DESC`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	queryWords := global.Parser.Cut(c.Query("info"), true)
	var tc58Jobs []model.TC58Job
	if err := global.Database.Select(&tc58Jobs, sqlString, strings.Join(queryWords, " | "),
		"%"+c.Query("area")+"%", strings.Join(queryWords, " | ")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []TC58JobResponse
	for _, tc58Job := range tc58Jobs {
		sqlString = `SELECT count(*) FROM user_favorite_58_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), tc58Job.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, TC58JobResponse{
			JobId:       tc58Job.ID,
			JobName:     tc58Job.JobName,
			JobArea:     tc58Job.JobArea,
			Salary:      tc58Job.Salary,
			JobWel:      strings.Split(tc58Job.JobWel, " "),
			CompanyName: tc58Job.CompanyName,
			JobNeed:     strings.Split(tc58Job.JobNeed, " "),
			JobURL:      tc58Job.JobURL,
			CreatedAt:   tc58Job.CreatedAt,
			IsFull:      tc58Job.IsFull,
			IsFavor:     isFavor == 1,
		})
	}
	c.JSON(http.StatusOK, All58Data{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

// Search58DataByRandom godoc
// @Schemes http
// @Description 随机搜索58同城数据
// @Tags Search
// @Param area query string false "地区"
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} All58Data "58同城数据"
// @Failure default {string} string "服务器错误"
// @Router /58_data/random [get]
// @Security ApiKeyAuth
func Search58DataByRandom(c *gin.Context) {
	sqlString := `SELECT * FROM "58_data" WHERE job_area LIKE $1 ORDER BY random()`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var tc58Jobs []model.TC58Job
	if err := global.Database.Select(&tc58Jobs, sqlString, "%"+c.Query("area")+"%"); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []TC58JobResponse
	for _, tc58Job := range tc58Jobs {
		sqlString = `SELECT count(*) FROM user_favorite_58_data WHERE user_id = $1 AND data_id = $2`
		var isFavor int
		if err := global.Database.Get(&isFavor, sqlString, c.GetInt("UserId"), tc58Job.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		jobs = append(jobs, TC58JobResponse{
			JobId:       tc58Job.ID,
			JobName:     tc58Job.JobName,
			JobArea:     tc58Job.JobArea,
			Salary:      tc58Job.Salary,
			JobWel:      strings.Split(tc58Job.JobWel, " "),
			CompanyName: tc58Job.CompanyName,
			JobNeed:     strings.Split(tc58Job.JobNeed, " "),
			JobURL:      tc58Job.JobURL,
			CreatedAt:   tc58Job.CreatedAt,
			IsFull:      tc58Job.IsFull,
			IsFavor:     isFavor == 1,
		})
	}
	c.JSON(http.StatusOK, All58Data{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}
