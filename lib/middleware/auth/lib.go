package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	Secretkey  []byte
	ExpireTime time.Duration
}

func NewJWT(key string, expiretime time.Duration) *JWT {
	return &JWT{
		Secretkey:  []byte(key),
		ExpireTime: expiretime,
	}
}

type Token_T struct {
	UID string `json:"userid"`
}

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func validateHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(data map[string]interface{}) (string, error) {
	date := time.Now().Add(time.Second * time.Duration(expireTime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   data,
		"exipre": date.Unix(),
	})

	tokenString, err := token.SignedString(secretkey)
	return tokenString, err
}
