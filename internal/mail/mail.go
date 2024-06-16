package mail

import (
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

func SendEmail(from, to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
