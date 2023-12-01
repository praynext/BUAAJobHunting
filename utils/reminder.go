package utils

import (
	"BUAAJobHunting/global"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"time"
)

type ReminderSend struct {
	ReminderId int    `db:"id"`
	Message    string `db:"message"`
	Email      string `db:"email"`
}

func SendEmailReminder(address string, message string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("BUAAJobHunting <BUAAJobHunting@163.com>")
	e.To = make([]string, 1)
	e.To[0] = address
	content := fmt.Sprintf(`
	<div>
		<div>
			尊敬的%s，您好！
		</div>
		<div style="padding: 8px 40px 8px 50px;">
			<p>%s</p>
		</div>
		<div>
			<p>此邮箱为系统邮箱，请勿回复。</p>
		</div>
	</div>
	`, address, message)
	e.HTML = []byte(content)
	err := e.SendWithTLS("smtp.163.com:465", global.SMTPAuth, &tls.Config{ServerName: "smtp.163.com"})
	return err
}

func CheckReminder() {
	sqlString := `SELECT r.id, r.message, u.email from reminder r join "user" u on u.id = r.user_id WHERE r.time < $1 AND r.has_sent = false`
	var reminderSends []ReminderSend
	if err := global.Database.Select(&reminderSends, sqlString, time.Now().In(time.FixedZone("CST", 8*3600))); err != nil {
		return
	}
	for _, reminderSend := range reminderSends {
		if err := SendEmailReminder(reminderSend.Email, reminderSend.Message); err != nil {
			continue
		}
		sqlString = `UPDATE reminder SET has_sent = true WHERE id = $1`
		if _, err := global.Database.Exec(sqlString, reminderSend.ReminderId); err != nil {
			continue
		}
	}
}
