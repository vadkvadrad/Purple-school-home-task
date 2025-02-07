package verification

import (
	"verify-api/internal/user"
)

type VerificationService struct {
	UserRepository *user.UserRepository
}

func NewVerificationService(userRepository *user.UserRepository) *VerificationService {
	return &VerificationService{
		UserRepository: userRepository,
	}
}