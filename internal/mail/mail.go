package mail

import (
	gomail "gopkg.in/mail.v2"

	log "github.com/sirupsen/logrus"
)

type Mail struct {
	SenderAddress string
	SmtpHost      string
	SmtpPort      int
	SmtpPassword  string
}

type MailBody struct {
	ReceiversAddress string
	Subject          string
	Body             string
}

func NewMail(mail Mail) *Mail {
	return &Mail{
		SenderAddress: mail.SenderAddress,
		SmtpHost:      mail.SmtpHost,
		SmtpPort:      mail.SmtpPort,
		SmtpPassword:  mail.SmtpPassword,
	}
}

func (a *Mail) Send(body MailBody) {
	m := gomail.NewMessage()
	m.SetHeader("From", a.SenderAddress)
	m.SetHeader("To", body.ReceiversAddress)
	m.SetHeader("Subject", body.Subject)
	m.SetBody("text/html", body.Body)

	d := gomail.NewDialer(a.SmtpHost, a.SmtpPort, a.SenderAddress, a.SmtpPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Panicln(err)
	}

	return
}
