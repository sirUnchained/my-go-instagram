package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey string
	aud       string
	iss       string
}

func (ja *JWTAuthenticator) GenerateToken(clams jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)

	token_str, err := token.SignedString([]byte(ja.secretKey))
	if err != nil {
		return "", err
	}

	return token_str, nil
}

func (ja *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token: %v", t.Header["alg"])
		}
		return []byte(ja.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(ja.aud),
		jwt.WithIssuer(ja.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
