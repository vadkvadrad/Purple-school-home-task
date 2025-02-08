package verify

import (
	"verify-api/internal/user"
)

type VerifyService struct {
	userRepository *user.UserRepository
}

func NewVerifyService(userRepository *user.UserRepository) *VerifyService {
	return &VerifyService{
		userRepository: userRepository,
	}
}

func(service * VerifyService) GetByEmail(email string) (*user.User, error) {
	return service.userRepository.FindByEmail(email)
}

func(service * VerifyService) GetByHash(hash string) (*user.User, error) {
	return service.userRepository.FindByHash(hash)
}

func(service * VerifyService) Create(email, password, address, hash string) (*user.User, error) {
	return service.userRepository.Create(&user.User{
		Email: email,
		Password: password,
		Address: address,
		Hash: hash,
		IsVerified: false,
	})
}

func(service * VerifyService) Verify(hash string) (*user.User, error) {
	return service.userRepository.Verify(hash)
}

func(service * VerifyService) GenerateHash() (string) {
	hash := user.GenerateHash()
	for {
		existedUser, _ := service.GetByHash(hash)
		if existedUser == nil {
			break
		}
		hash = user.GenerateHash()
	}
	return hash
}