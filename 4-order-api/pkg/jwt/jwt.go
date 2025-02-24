package jwt

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Phone string
}

type JWT struct {
	Secret string
}

func NewJwt(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (token string, err error) {
	bytes := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": data.Phone,
		"exp": json.Number(strconv.FormatInt(time.Now().Add(time.Hour*time.Duration(1)).Unix(), 10)),
		"iat": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	})

	token, err = bytes.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}


func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}

	phone := t.Claims.(jwt.MapClaims)["phone"]

	return true, &JWTData{
		Phone: phone.(string),
	}
} 