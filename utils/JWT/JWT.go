package JWT

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userid int, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	exp := now.UTC().Add(ttl).Unix()
	claims := jwt.MapClaims{
		"id_user": userid,
		"exp":     exp,
		"iat":     now.Unix(),
	}

	JWTtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JWT, err := JWTtoken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return JWT, nil
}

func ValidateToken(token string, secret string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Ложный метод подписи, метод - %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return tok.Claims.(jwt.MapClaims), err
}
