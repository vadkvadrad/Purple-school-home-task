package verify

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
	"verify-api/configs"
	"verify-api/internal/user"
	"verify-api/pkg/er"

	"github.com/jordan-wright/email"
)

type VerifyService struct {
	UserRepository *user.UserRepository
	Config *configs.Config
}

type VerifyServiceDeps struct {
	UserRepository *user.UserRepository
	Config *configs.Config
}

func NewVerifyService(deps VerifyServiceDeps) *VerifyService {
	return &VerifyService{
		UserRepository: deps.UserRepository,
		Config: deps.Config,
	}
}

func(service *VerifyService) GetByEmail(email string) (*user.User, error) {
	return service.UserRepository.FindByEmail(email)
}

func(service *VerifyService) GetByHash(hash string) (*user.User, error) {
	return service.UserRepository.FindByHash(hash)
}

func(service *VerifyService) Create(email, hash string) (*user.User, error) {
	return service.UserRepository.Create(&user.User{
		Email: email,
		Hash: hash,
		IsVerified: false,
	})
}

func(service *VerifyService) Verify(hash string) (*user.User, error) {
	return service.UserRepository.Verify(hash)
}

func(service *VerifyService) GenerateHash() (string) {
	hash := user.GenerateHash()
	for {
		existedUser, _ := service.UserRepository.FindByHash(hash)
		if existedUser == nil {
			break
		}
		hash = user.GenerateHash()
	}
	return hash
}

func (service *VerifyService) Send(hash string) error {
	user, err := service.UserRepository.FindByHash(hash)
	if err != nil {
		return err
	}

	// Настроим письмо
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", service.Config.Sender.Name, service.Config.Sender.Email)
	e.To = []string{user.Email}
	e.Subject = "Подтверждение электронной почты"
	e.Text = []byte("Для подтверждения электронной почты перейдите по ссылке: http://localhost:8081/verify/" + user.Hash)

	// Настройки SMTP
	server := service.Config.Sender.Address
	port := service.Config.Sender.Port
	auth := smtp.PlainAuth("", service.Config.Sender.Email, service.Config.Sender.Password, server)

	// Настроим таймаут подключения
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
	}

	// Настроим TLS
	tlsConfig := &tls.Config{
		ServerName: server,
	}

	// Подключаемся к серверу
	conn, err := tls.DialWithDialer(dialer, "tcp", server+":"+port, tlsConfig)
	if err != nil {
		return er.Wrap("Ошибка подключения:", err)
	}

	// Создаём SMTP-клиента
	c, err := smtp.NewClient(conn, server)
	if err != nil {
		return er.Wrap("Ошибка SMTP-клиента:", err)
	}

	// Аутентификация
	if err = c.Auth(auth); err != nil {
		return er.Wrap("Ошибка аутентификации:", err)
	}

	// Отправляем письмо
	if err = e.SendWithTLS(server+":"+port, auth, tlsConfig); err != nil {
		return er.Wrap("Ошибка отправки письма:", err)
	}

	return nil
}