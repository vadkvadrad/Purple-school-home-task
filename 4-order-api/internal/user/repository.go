package user

import (
	"fmt"
	"order-api/pkg/db"

	"gorm.io/gorm/clause"
)

const (
	EmailKey     = "email"
	PhoneKey     = "phone"
	SessionIdKey = "session_id"
)

type UserRepository struct {
	Db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByKey(key, data string) (*User, error) {
	var user User
	query := fmt.Sprintf("%s = ?", key)
	result := repo.Db.First(&user, query, data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
