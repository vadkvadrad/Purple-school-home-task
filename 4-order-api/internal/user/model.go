package user

import (
	"math/rand"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	SessionId string `json:"session_id" gorm:"uniqueIndex"`
	Code string `json:"code"`
}


func (u *User) Generate() {
	u.SessionId = randLettersRunes(10)
	u.Code = randNumbersRunes(4)
}

var lettersRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbersRunes = []rune("0123456789")

func randLettersRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}

func randNumbersRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = numbersRunes[rand.Intn(len(numbersRunes))]
	}
	return string(b)
}