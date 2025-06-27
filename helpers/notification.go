package helpers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Notification struct {
	From     string
	To       string
	Password string
	Host     string
	Port     string
}

func NewNotification(from, to, password, host, port string) (*Notification, error) {
	if from == "" || to == "" || password == "" || host == "" || port == "" {
		return nil, errors.New("all fields are required")
	}
	return &Notification{
		From:     from,
		To:       to,
		Password: password,
		Host:     host,
		Port:     port,
	}, nil
}

func (n *Notification) SendEmail(subject, body string) error {
	if subject == "" || body == "" {
		return errors.New("subject and body are required")
	}
	if !strings.Contains(subject, "=") && !strings.Contains(body, "=") {
		msg := "To: " + n.To + "\r\nSubject: " + subject + "\r\n\r\n" + body
		auth := smtp.PlainAuth("", n.From, n.Password, n.Host)
		err := smtp.SendMail(n.Host+":"+n.Port, auth, n.From, []string{n.To}, []byte(msg))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("potential security vulnerability detected")
}

func (n *Notification) SendSuccessNotification(jobName string) error {
	if jobName == "" {
		return errors.New("job name is required")
	}
	subject := "Job Completed Successfully"
	body := "The job " + jobName + " has been completed successfully."
	return n.SendEmail(subject, body)
}

func (n *Notification) SendErrorNotification(jobName, errorMessage string) error {
	if jobName == "" || errorMessage == "" {
		return errors.New("job name and error message are required")
	}
	subject := "Job Failed"
	body := "The job " + jobName + " has failed with the following error: " + errorMessage
	return n.SendEmail(subject, body)
}