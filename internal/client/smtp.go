package client

import (
	"currency-checker/internal/config"
	"gopkg.in/gomail.v2"
)

type SmtpClient struct {
	AuthorEmail string
	Password    string
	Host        string
	Port        int
}

func NewSmtpClient(cfg *config.Config) *SmtpClient {
	return &SmtpClient{
		AuthorEmail: cfg.EmailServer.Credentials.Email,
		Password:    cfg.EmailServer.Credentials.Password,
		Host:        cfg.EmailServer.Host,
		Port:        cfg.EmailServer.Port,
	}
}

func (s *SmtpClient) SendEmail(email string, priceInfo string) error {

	message := gomail.NewMessage()
	message.SetHeader("From", s.AuthorEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "BTC Price")
	message.SetBody("text/plain", priceInfo)

	dialer := gomail.NewDialer(s.Host, s.Port, s.AuthorEmail, s.Password)

	err := dialer.DialAndSend(message)
	if err != nil {
		return err
	}
	return nil
}
