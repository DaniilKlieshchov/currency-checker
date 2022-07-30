package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type EmailService struct {
	storage        Storage
	coinbaseClient Client
	smtpClient     EmailSender
}

func NewEmailService(storage Storage, client Client, smtpClient EmailSender) *EmailService {
	return &EmailService{storage: storage, coinbaseClient: client, smtpClient: smtpClient}
}

type Storage interface {
	Append(email string) error
	GetEmails() []string
}

type Client interface {
	GetRate() (int, error)
}

type EmailSender interface {
	SendEmail(string, string) error
}

func (s *EmailService) Rate() (int, error) {
	rate, err := s.coinbaseClient.GetRate()
	if err != nil {
		return 0, err
	}
	return rate, nil
}

func (s *EmailService) SendEmails(log logrus.FieldLogger) []string {
	price, err := s.coinbaseClient.GetRate()
	info := fmt.Sprintf("Current BTC price: %d UAH", price)
	if err != nil {
		log.Error(err)
		return nil
	}
	emails := s.storage.GetEmails()
	mutex := sync.Mutex{}
	failedEmails := make([]string, 0, len(emails))
	wg := sync.WaitGroup{}
	wg.Add(len(emails))
	for _, email := range emails {
		email := email
		go func() {
			defer wg.Done()
			err := s.smtpClient.SendEmail(email, info)
			if err != nil {
				mutex.Lock()
				failedEmails = append(failedEmails, email)
				mutex.Unlock()
				log.Error(err.Error())
			}
		}()
	}
	wg.Wait()
	return failedEmails
}

func (s *EmailService) Subscribe(email string) error {
	err := s.storage.Append(email)
	if err != nil {
		return err
	}
	return nil
}
