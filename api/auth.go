package api

import (
	"BUAAJobHunting/global"
	"BUAAJobHunting/model"
	"BUAAJobHunting/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

const (
	lastSentTimesKey = "last_sent_times"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name     string  `json:"name" bind:"required,max=20"`
	Password string  `json:"password" bind:"required,min=6,max=20"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	VCode    string  `json:"v_code"`
}

type ChangePwdRequest struct {
	OldPassword string `json:"old_password" bind:"required"`
	NewPassword string `json:"new_password" bind:"required,min=6,max=20"`
}

type ResetPwdRequest struct {
	UserName    string `json:"username" bind:"required"`
	VCode       string `json:"v_code" bind:"required"`
	NewPassword string `json:"new_password" bind:"required,min=6,max=20"`
}

// Login godoc
// @Schemes http
// @Description 用户登录
// @Tags Authentication
// @Param info body LoginRequest true "用户登陆信息"
// @Success 200 {object} LoginResponse 用户登陆反馈
// @Failure 400 {string} string "请求解析失败"/"密码错误"
// @Failure 404 {string} string "用户名不存在"
// @Failure default {string} string "服务器错误"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	var userInfo model.User
	sqlString := `SELECT id, password FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, loginRequest.UserName); err != nil {
		c.String(http.StatusNotFound, "用户名不存在")
		return
	}
	if !utils.VerifyPassword(userInfo.Password, loginRequest.Password) {
		c.String(http.StatusBadRequest, "密码错误")
		return
	}
	token, err := global.CreateSession(c, &global.Session{
		Role:   global.USER,
		UserId: userInfo.ID,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
	c.Set("Role", global.USER)
	c.Set("UserId", userInfo.ID)
}

// Register godoc
// @Schemes http
// @Description 用户注册
// @Tags Authentication
// @Param info body RegisterRequest true "用户注册信息"
// @Success 200 {string} string "注册成功"
// @Failure 400 {string} string "请求解析失败"/"验证码已过期"/"验证码错误"
// @Failure 409 {string} string "用户名已存在"
// @Failure default {string} string "服务器错误"
// @Router /register [post]
func Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	var userInfo model.User
	sqlString := `SELECT id FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, registerRequest.Name); err == nil {
		c.String(http.StatusConflict, "用户名已存在")
		return
	}
	rawCode := global.Redis.Get(c, registerRequest.Email)
	if rawCode.Err() != nil {
		c.String(http.StatusBadRequest, "验证码已过期")
		return
	} else if rawCode.Val() != registerRequest.VCode {
		c.String(http.StatusBadRequest, "验证码错误")
		return
	} else {
		global.Redis.Del(c, registerRequest.Email)
	}
	userInfo.Name = registerRequest.Name
	var err error
	userInfo.Password, err = utils.EncryptPassword(registerRequest.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	userInfo.Email = registerRequest.Email
	userInfo.Phone = registerRequest.Phone
	sqlString = `INSERT INTO "user" (name, password, email, phone, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := global.Database.Get(&userInfo.ID, sqlString, userInfo.Name, userInfo.Password,
		userInfo.Email, userInfo.Phone, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}

// ChangePassword godoc
// @Schemes http
// @Description 用户修改密码
// @Tags Authentication
// @Param info body ChangePwdRequest true "用户修改密码信息"
// @Success 200 {string} string "修改成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 409 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /change_password [post]
// @Security ApiKeyAuth
func ChangePassword(c *gin.Context) {
	var changePwdRequest ChangePwdRequest
	if err := c.ShouldBindJSON(&changePwdRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userId := c.GetInt("UserId")
	var userInfo model.User
	sqlString := `SELECT id, password FROM "user" WHERE id = $1`
	if err := global.Database.Get(&userInfo, sqlString, userId); err != nil {
		c.String(http.StatusConflict, "修改密码失败")
		return
	}
	if !utils.VerifyPassword(userInfo.Password, changePwdRequest.OldPassword) {
		c.String(http.StatusBadRequest, "修改密码失败")
		return
	}
	var err error
	userInfo.Password, err = utils.EncryptPassword(changePwdRequest.NewPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE "user" SET password = $1 WHERE id = $2`
	if _, err := global.Database.Exec(sqlString, userInfo.Password, userInfo.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "修改成功")
}

// ResetPassword godoc
// @Schemes http
// @Description 用户重置密码，需要先向邮箱发送一封邮件
// @Tags Authentication
// @Param info body ResetPwdRequest true "用户修改密码信息"
// @Success 200 {string} string "修改成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 404 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /reset_password [post]
// @Security ApiKeyAuth
func ResetPassword(c *gin.Context) {
	// 验证Redis内的邮箱验证码是否正确，然后修改密码
	var resetPwdRequest ResetPwdRequest
	if err := c.ShouldBindJSON(&resetPwdRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	var userInfo model.User
	sqlString := `SELECT * FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, resetPwdRequest.UserName); err != nil {
		c.String(http.StatusNotFound, "修改密码失败")
		return
	}
	rawCode := global.Redis.Get(c, userInfo.Email)
	if rawCode.Err() != nil {
		c.String(http.StatusBadRequest, "修改密码失败")
		return
	} else if rawCode.Val() != resetPwdRequest.VCode {
		c.String(http.StatusBadRequest, "修改密码失败")
		return
	} else {
		global.Redis.Del(c, userInfo.Email)
	}
	var err error
	userInfo.Password, err = utils.EncryptPassword(resetPwdRequest.NewPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE "user" SET password = $1 WHERE id = $2`
	if _, err := global.Database.Exec(sqlString, userInfo.Password, userInfo.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "修改成功")
}

// Logout godoc
// @Schemes http
// @Description 用户退出
// @Tags Authentication
// @Success 200 {string} string "退出成功"
// @Failure default {string} string "服务器错误"
// @Router /logout [post]
// @Security ApiKeyAuth
func Logout(c *gin.Context) {
	err := global.DeleteSession(c, c.Request.Header.Get(global.TokenHeader))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
	}
	c.String(http.StatusOK, "退出成功")
}

// SendEmail godoc
// @Schemes http
// @Description 发送邮件
// @Tags Authentication
// @Param email query string true "邮箱"
// @Success 200 {string} string "发送成功"
// @Failure 400 {string} string "验证码存储失败"
// @Failure default {string} string "服务器错误"
// @Router /send_email [post]
func SendEmail(c *gin.Context) {
	email := c.Query("email")
	currentTime := time.Now()
	lastSentTimeStr, err := global.Redis.ZScore(c, lastSentTimesKey, email).Result()
	if err == redis.Nil || (err == nil && currentTime.Sub(time.Unix(int64(lastSentTimeStr), 0)) >= time.Minute) {
		vCode, err := utils.SendEmailValidate(c.Query("email"))
		if err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		err = global.Redis.ZAdd(c, lastSentTimesKey, redis.Z{
			Score:  float64(currentTime.Unix()),
			Member: email,
		}).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "验证码存储失败")
			return
		}
		err = global.Redis.Set(c, email, vCode, time.Minute*5).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "验证码存储失败")
			return
		}
		c.String(http.StatusOK, "发送成功")
	} else {
		c.String(http.StatusBadRequest, "发送过于频繁")
	}
}

type UserResponse struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type AllUserData struct {
	TotalCount int            `json:"total_count"`
	Users      []UserResponse `json:"users"`
}

// GetAllUser godoc
// @Schemes http
// @Description 获取所有用户信息
// @Tags Authentication
// @Success 200 {object} AllUserData "所有用户信息"
// @Failure default {string} string "服务器错误"
// @Router /all_user [get]
// @Security ApiKeyAuth
func GetAllUser(c *gin.Context) {
	sqlString := `SELECT id, name FROM "user"`
	var users []model.User
	if err := global.Database.Select(&users, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			UserId:   user.ID,
			UserName: user.Name,
		})
	}
	c.JSON(http.StatusOK, AllUserData{
		TotalCount: len(users),
		Users:      userResponses,
	})
}
