package sendmail

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type SendMail struct {
	username string
	password string
	domain   string
	subject  string
	message  string
}

func NewSendMail(username, password, domain, subject, message string) *SendMail {
	return &SendMail{
		username: username,
		password: password,
		domain:   domain,
		subject:  subject,
		message:  message,
	}
}

func (this *SendMail) Send(to string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", this.username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", this.subject)
	m.SetBody("text/html", fmt.Sprintf(`<b>%s</b> %s`, message, this.message))

	d := gomail.NewPlainDialer(this.domain, 587, this.username, this.password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
