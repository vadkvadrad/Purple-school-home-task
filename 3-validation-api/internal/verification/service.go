package verification

import "verify-api/pkg/db"

type VerificationService struct {
	Db *db.Db
}

func NewVerificationService(db *db.Db) *VerificationService {
	return &VerificationService{
		Db: db,
	}
}