package auth

import (
	"fmt"
	"order-api/internal/user"
)


type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
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

		service.UserRepository.Create(&createdUser)

		err := service.send(&createdUser)
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
	fmt.Println(user.Code)
	return nil
}