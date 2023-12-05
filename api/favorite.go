package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// UserFavorite58Data godoc
// @Schemes http
// @Description 用户收藏58同城数据
// @Tags Favorite
// @Param id query int true "58同城数据ID"
// @Success 200 {string} string "收藏成功"
// @Failure default {string} string "服务器错误"
// @Router /58_data/favorite [post]
// @Security ApiKeyAuth
func UserFavorite58Data(c *gin.Context) {
	sqlString := `INSERT INTO user_favorite_58_data (user_id, data_id, created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Query("id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "收藏成功")
}

// UserFavoriteBossData godoc
// @Schemes http
// @Description 用户收藏Boss直聘数据
// @Tags Favorite
// @Param id query int true "Boss直聘数据ID"
// @Success 200 {string} string "收藏成功"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/favorite [post]
// @Security ApiKeyAuth
func UserFavoriteBossData(c *gin.Context) {
	sqlString := `INSERT INTO user_favorite_boss_data (user_id, data_id, created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Query("id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "收藏成功")
}

// UserCancelFavorite58Data godoc
// @Schemes http
// @Description 用户取消收藏58同城数据
// @Tags Favorite
// @Param id query int true "58同城数据ID"
// @Success 200 {string} string "取消收藏成功"
// @Failure default {string} string "服务器错误"
// @Router /58_data/favorite [delete]
// @Security ApiKeyAuth
func UserCancelFavorite58Data(c *gin.Context) {
	sqlString := `DELETE FROM user_favorite_58_data WHERE user_id=$1 AND data_id=$2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Query("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消收藏成功")
}

// UserCancelFavoriteBossData godoc
// @Schemes http
// @Description 用户取消收藏Boss直聘数据
// @Tags Favorite
// @Param id query int true "Boss直聘数据ID"
// @Success 200 {string} string "取消收藏成功"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/favorite [delete]
// @Security ApiKeyAuth
func UserCancelFavoriteBossData(c *gin.Context) {
	sqlString := `DELETE FROM user_favorite_58_data WHERE user_id=$1 AND data_id=$2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Query("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消收藏成功")
}

// UserGetFavorite58Data godoc
// @Schemes http
// @Description 用户获取收藏的58同城数据
// @Tags Favorite
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} All58Data "58同城数据"
// @Failure default {string} string "服务器错误"
// @Router /58_data/favorite [get]
// @Security ApiKeyAuth
func UserGetFavorite58Data(c *gin.Context) {
	sqlString := `
		SELECT * FROM "58_data" 
		WHERE id IN (
			SELECT data_id FROM "user_favorite_58_data" WHERE user_id=$1
		)
	`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var tc58Jobs []model.TC58Job
	if err := global.Database.Select(&tc58Jobs, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []TC58JobResponse
	for _, tc58Job := range tc58Jobs {
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
		})
	}
	c.JSON(http.StatusOK, All58Data{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}

// UserGetFavoriteBossData godoc
// @Schemes http
// @Description 用户获取收藏的Boss直聘数据
// @Tags Favorite
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllBossData "Boss直聘数据"
// @Failure default {string} string "服务器错误"
// @Router /boss_data/favorite [get]
// @Security ApiKeyAuth
func UserGetFavoriteBossData(c *gin.Context) {
	sqlString := `
		SELECT * FROM "boss_data" 
		WHERE id IN (
			SELECT data_id FROM "user_favorite_boss_data" WHERE user_id=$1
		)
	`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var bossJobs []model.BossJob
	if err := global.Database.Select(&bossJobs, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var jobs []BossJobResponse
	for _, bossJob := range bossJobs {
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
		})
	}
	c.JSON(http.StatusOK, AllBossData{
		TotalCount: len(jobs),
		Jobs:       jobs,
	})
}
