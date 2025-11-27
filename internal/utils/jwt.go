package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

type JwtClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

func InitJwt(secret string){
	jwtSecret = []byte(secret)
}

func GenerateJwt(userid string, expirationHours int) (string, error){
	claims := jwt.MapClaims{
		"sub":userid,
		"iat":time.Now().Unix(),
		"exp":time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (string, error){
	if jwtSecret == nil {
		return "", errors.New("jwt not initiased")
	}

	claims := &JwtClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		method, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != jwt.SigningMethodHS256{
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Sub == ""{
		return "", errors.New("missing token subject")
	}

	return claims.Sub, nil
}