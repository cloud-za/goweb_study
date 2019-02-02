package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

var secretkey []byte
var expireTime int = 7 * 24 * 60 * 60

func Setup(key []byte, ex int) {
	secretkey = key
	expireTime = ex
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
