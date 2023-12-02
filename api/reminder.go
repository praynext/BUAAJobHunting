package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ReminderResponse struct {
	ReminderId int       `json:"reminder_id"`
	UserId     int       `json:"user_id"`
	Message    string    `json:"message"`
	Time       string    `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
	HasSent    bool      `json:"has_sent"`
}

type AllReminderData struct {
	TotalCount int                `json:"total_count"`
	Reminders  []ReminderResponse `json:"reminders"`
}

// GetReminder godoc
// @Schemes http
// @Description 用户获取提醒事项
// @Tags Reminder
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllReminderData "提醒事项"
// @Failure default {string} string "服务器错误"
// @Router /reminder/get [get]
// @Security ApiKeyAuth
func GetReminder(c *gin.Context) {
	sqlString := `SELECT * from reminder WHERE user_id=$1`
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	var reminders []model.Reminder
	if err := global.Database.Select(&reminders, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var reminderResponses []ReminderResponse
	for _, reminder := range reminders {
		reminderResponses = append(reminderResponses, ReminderResponse{
			ReminderId: reminder.ID,
			UserId:     reminder.UserId,
			Message:    reminder.Message,
			Time:       reminder.Time.Format("2006/01/02 15:04"),
			CreatedAt:  reminder.CreatedAt,
			HasSent:    reminder.HasSent,
		})
	}
	c.JSON(http.StatusOK, AllReminderData{
		TotalCount: len(reminderResponses),
		Reminders:  reminderResponses,
	})
}

type AddReminderRequest struct {
	Message string `json:"message" binding:"required"`
	Time    string `json:"time" binding:"required"`
}

// AddReminder godoc
// @Schemes http
// @Description 用户添加提醒事项
// @Tags Reminder
// @Param info body AddReminderRequest true "提醒事项信息"
// @Success 200 {string} string "添加提醒事项成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /reminder/add [post]
// @Security ApiKeyAuth
func AddReminder(c *gin.Context) {
	var addReminderRequest AddReminderRequest
	if err := c.ShouldBindJSON(&addReminderRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	formatTime, err := time.Parse("2006/01/02 15:04", addReminderRequest.Time)
	if err != nil {
		c.String(http.StatusBadRequest, "时间格式错误")
		return
	}
	sqlString := `INSERT INTO reminder (user_id, message, time, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"),
		addReminderRequest.Message, formatTime, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加提醒事项成功")
}

type UpdateReminderRequest struct {
	ReminderId int    `json:"reminder_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
	Time       string `json:"time" binding:"required"`
}

// UpdateReminder godoc
// @Schemes http
// @Description 用户更新提醒事项
// @Tags Reminder
// @Param info body UpdateReminderRequest true "提醒事项信息"
// @Success 200 {string} string "更新提醒事项成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /reminder/update [put]
// @Security ApiKeyAuth
func UpdateReminder(c *gin.Context) {
	var updateReminderRequest UpdateReminderRequest
	if err := c.ShouldBindJSON(&updateReminderRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	formatTime, err := time.Parse("2006/01/02 15:04", updateReminderRequest.Time)
	if err != nil {
		c.String(http.StatusBadRequest, "时间格式错误")
		return
	}
	sqlString := `UPDATE reminder SET message=$1, time=$2 WHERE id=$3 AND user_id=$4`
	if _, err := global.Database.Exec(sqlString, updateReminderRequest.Message, formatTime,
		updateReminderRequest.ReminderId, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新提醒事项成功")
}

// DeleteReminder godoc
// @Schemes http
// @Description 用户删除提醒事项
// @Tags Reminder
// @Param id query int true "提醒事项ID"
// @Success 200 {string} string "删除提醒事项成功"
// @Failure default {string} string "服务器错误"
// @Router /reminder/delete [delete]
// @Security ApiKeyAuth
func DeleteReminder(c *gin.Context) {
	sqlString := `DELETE FROM reminder WHERE id=$1 AND user_id=$2`
	if _, err := global.Database.Exec(sqlString, c.Query("id"), c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除提醒事项成功")
}
