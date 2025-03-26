package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

func DecodeWithouSignature(token string) (jwt.Claims, error) {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token, &claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
