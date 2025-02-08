package user

import (
	"math/rand"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
	IsVerified bool `json:"is_verified"`
}

func GenerateHash() string {
	return RandStringRunes(6)
}

var lettersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}