package user

import "verify-api/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, result.Error
}

func (repo *UserRepository) FindByHash(hash string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) Verify(hash string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	user.IsVerified = true
	result = repo.Database.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}