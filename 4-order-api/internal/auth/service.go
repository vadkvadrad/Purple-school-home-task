package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"order-api/configs"
	"order-api/internal/user"
	"order-api/pkg/er"
	"time"

	"github.com/jordan-wright/email"
)


type AuthService struct {
	UserRepository *user.UserRepository
	Config *configs.Config
}

func NewAuthService(config *configs.Config, userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		Config: config,
	}
}

func (service *AuthService) GetBySessionId(sessionId string) (*user.User, error) {
	existedUser, err := service.UserRepository.FindByKey(user.SessionIdKey, sessionId)
	if err != nil {
		return nil, err
	}
	return existedUser, nil
}

func (service *AuthService) Login(phone, email string) (string, error) {
	existedUser, _ := service.UserRepository.FindByKey(user.PhoneKey, phone)
	if existedUser == nil {
		createdUser := user.User{
			Phone: phone,
			Email: email,
		}
		createdUser.Generate()

		_, err := service.UserRepository.Create(&createdUser)
		if err != nil {
			return "", err
		}

		err = service.send(&createdUser)
		if err != nil {return "", err}
		return createdUser.SessionId, nil
	}
	existedUser.Email = email
	existedUser.Generate()

	service.UserRepository.Update(existedUser)

	err := service.send(existedUser)
	if err != nil {
		return "", err
	}
	return existedUser.SessionId, nil
}

func (service *AuthService) send(user *user.User) error {
	// here will be sms registration
	fmt.Println(user.Code)

	if user.Email != "" {
		err := service.sendOnEmail(user.Email, user.Code)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *AuthService) sendOnEmail(emailAddr string, code string) error {
	// Настроим письмо
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", service.Config.Sender.Name, service.Config.Sender.Email)
	e.To = []string{emailAddr}
	e.Subject = "Код подтверждения личности"
	e.Text = []byte("Ваш персональный код подтверждения личности: " + code + ". Не сообщайте никому данный код.")

	// Настройки SMTP
	server := service.Config.Sender.Address
	port := service.Config.Sender.Port
	err := validate(server, port)
	if err != nil {
		return err
	}
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
	defer c.Close()

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

func validate(server string, port string) error {
	if server == "" {
		return errors.New("server is not specified")
	}
	if port == "" {
		return errors.New("port is not specified")
	}
	return nil
}