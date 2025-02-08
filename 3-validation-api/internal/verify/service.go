package verify

import (
	"verify-api/internal/user"
)

type VerifyService struct {
	UserRepository *user.UserRepository
}

func NewVerifyService(userRepository *user.UserRepository) *VerifyService {
	return &VerifyService{
		UserRepository: userRepository,
	}
}