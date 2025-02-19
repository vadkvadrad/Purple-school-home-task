package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	SessionId string `json:"session_id" gorm:"uniqueIndex"`
	Code string `json:"code"`
}
