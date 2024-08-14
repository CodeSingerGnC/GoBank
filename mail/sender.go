package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAdress = "smtp.gmail.com"
	smtpServerAderss = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type SinaSender struct {
	name string
	fromEmailAdress string
	fromEmailPassword string
}

func NewSinaSender(name, fromEmailAdress, fromEmailPassword string) EmailSender {
	return &SinaSender{
		name: name,
		fromEmailAdress: fromEmailAdress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *SinaSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAdress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc 

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAdress, sender.fromEmailPassword, smtpAuthAdress)
	return e.Send(smtpServerAderss, smtpAuth)
}