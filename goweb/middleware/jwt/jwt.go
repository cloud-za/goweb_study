package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	Secretkey      []byte
	ExpireTime     time.Duration
	HalfOfExpire   int64
	ExpireInSecond int
}

func NewJWT(key string, expiretime int) *JWT {
	expireDuration := time.Second * time.Duration(expiretime)
	return &JWT{
		Secretkey:      []byte(key),
		ExpireTime:     expireDuration,
		HalfOfExpire:   int64(expireDuration) / 2,
		ExpireInSecond: expiretime,
	}
}

//VerifyToken is a Middleware for Gin
func (j *JWT) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authhdr := c.Request.Header.Get("Authorization")
		if authhdr == "" {
			//Authorization header is not exist, check cookie
			t, err := c.Request.Cookie("access_token")
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
			} else {
				token = t.Value
			}
		} else if strings.Contains(authhdr, "Bearer ") {
			token = strings.TrimLeft(authhdr, "Bearer ")
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"status": -1,
				"msg":    "invalid token format",
			})
			return
		}

		//verify token
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"status": -1,
				"msg":    "invalid token",
			})
			return
		}

		//parse token
		claims, err := j.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		//validate expire time
		diff := claims.StandardClaims.ExpiresAt - time.Now().Unix()
		if diff < 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"status": -1,
				"msg":    "token expired",
			})
			return
		}
		/*
			if diff < j.HalfOfExpire {
				tokenNew, err := j.RefreshToken(token)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    fmt.Errorf("Refresh token error: %+v", err),
					})
					c.Abort()
					return
				}
				c.SetCookie("access_token", tokenNew, j.ExpireInSecond, "/", "", true, true)
			}*/
		c.Set("claims", claims)
	}
}

type TokenClaims struct {
	ID        string `json:"userId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Privilege int32  `json:"priv"`
	jwt.StandardClaims
}

func NewTokenClaims(id string, name string, phone string, priv int32) *TokenClaims {
	return &TokenClaims{
		ID:        id,
		Name:      name,
		Phone:     phone,
		Privilege: priv,
	}

}

//GenerateToken is used to generate token
func (j *JWT) GenerateToken(claims *TokenClaims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Second * time.Duration(j.ExpireTime)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secretkey)
}

func (j *JWT) keyLookup(token *jwt.Token) (interface{}, error) {
	return j.Secretkey, nil
}

// ParseToken is used to parse token
func (j *JWT) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, j.keyLookup)
	if err != nil {
		return nil, fmt.Errorf("token error: %+v", err)
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, j.keyLookup)
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Second * time.Duration(j.ExpireTime)).Unix()
		return j.GenerateToken(claims)
	}
	return "", errors.New("token is invalid")
}
