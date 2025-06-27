package helpers

import (
	"encoding/json"
	"errors"
	"log"
	"net/smtp"
)

type WpsService struct {
	smtpServer string
	smtpPort   string
	fromEmail  string
	password   string
}

func NewWpsService(smtpServer, smtpPort, fromEmail, password string) *WpsService {
	return &WpsService{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		fromEmail:  fromEmail,
		password:   password,
	}
}

func (w *WpsService) _notifyDownloadViaEmail(toEmail, subject, body string) error {
	if toEmail == "" || subject == "" || body == "" {
		return errors.New("invalid input parameters")
	}

	msg := "To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body

	auth := smtp.PlainAuth("", w.fromEmail, w.password, w.smtpServer)

	err := smtp.SendMail(w.smtpServer+":"+w.smtpPort, auth, w.fromEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (w *WpsService) _notifyErrorViaEmail(toEmail, subject, body, errorMessage string) error {
	if toEmail == "" || subject == "" || body == "" || errorMessage == "" {
		return errors.New("invalid input parameters")
	}

	msg := "To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\nError Message: " + errorMessage

	auth := smtp.PlainAuth("", w.fromEmail, w.password, w.smtpServer)

	err := smtp.SendMail(w.smtpServer+":"+w.smtpPort, auth, w.fromEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (w *WpsService) NotifyDownloadViaEmail(toEmail, subject, body string) {
	err := w._notifyDownloadViaEmail(toEmail, subject, body)
	if err != nil {
		log.Println(err)
	}
}

func (w *WpsService) NotifyErrorViaEmail(toEmail, subject, body, errorMessage string) {
	err := w._notifyErrorViaEmail(toEmail, subject, body, errorMessage)
	if err != nil {
		log.Println(err)
	}
}