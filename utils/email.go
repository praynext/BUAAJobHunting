package utils

import (
	"BUAAJobHunting/global"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"time"
)

func GenerateDigitalCode(n int) string {
	letters := []byte("0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

func SendEmailValidate(recipient string) (string, error) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("BUAAJobHunting <BUAAJobHunting@163.com>")
	e.To = make([]string, 1)
	e.To[0] = recipient
	vCode := GenerateDigitalCode(6)
	content := fmt.Sprintf(`
	<div>
		<div>
			尊敬的%s，您好！
		</div>
		<div style="padding: 8px 40px 8px 50px;">
			<p>您于 %s 提交的邮箱验证，本次验证码为<u><strong>%s</strong></u>，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
		</div>
		<div>
			<p>此邮箱为系统邮箱，请勿回复。</p>
		</div>
	</div>
	`, recipient, time.Now().Format("2006-01-02 15:04:05"), vCode)
	e.HTML = []byte(content)
	err := e.SendWithTLS("smtp.163.com:465", global.SMTPAuth, &tls.Config{ServerName: "smtp.163.com"})
	return vCode, err
}
