package auth

import (
	"fmt"
	"order-api/configs"
	"order-api/internal/user"
	"order-api/pkg/sender"
)


type AuthService struct {
	UserRepository *user.UserRepository
	Config *configs.Config
	Sender *sender.Sender
}

type AuthServiceDeps struct {
	UserRepository *user.UserRepository
	Config *configs.Config
	Sender *sender.Sender
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		UserRepository: deps.UserRepository,
		Config: deps.Config,
		Sender: deps.Sender,
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
		subject := "Код подтверждения личности"
		text := "Ваш персональный код подтверждения личности: " + user.Code + ". Не сообщайте никому данный код."
		err := service.Sender.Email(user.Email, subject, text)
		if err != nil {
			return err
		}
	}
	return nil
}