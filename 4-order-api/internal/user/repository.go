package user

import (
	"fmt"
	"order-api/configs"
	"order-api/pkg/db"
)

const (
	EmailKey = "email"
	PhoneKey = "phone"
	SessionIdKey = "session_id"
)

type UserRepository struct {
	Db *db.Db
	Config *configs.Config
}

func NewUserRepository(db *db.Db, config *configs.Config) *UserRepository {
	return &UserRepository{
		Db: db,
		Config: config,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Db.Create(user)
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