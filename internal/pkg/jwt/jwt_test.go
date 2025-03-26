package jwt_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	internalJwt "main-api/internal/pkg/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtDecodedWithoutSignature(t *testing.T) {
	t.Parallel()

	t.Run("should return claims without signature", func(t *testing.T) {
		token := generateFakeJwtToken()

		claims, err := internalJwt.DecodeWithouSignature(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		token := "invalid-token"

		claims, err := internalJwt.DecodeWithouSignature(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func generateFakeJwtToken() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating ECDSA key: %v\n", err)
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  1516239022,
	})

	stringToken, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return ""
	}

	return stringToken
}
