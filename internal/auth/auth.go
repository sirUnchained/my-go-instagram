package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(clams jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

func NewJWTAuthenticator(secretKey, aud, iss string) Authenticator {
	return &JWTAuthenticator{
		secretKey: secretKey,
		aud:       aud,
		iss:       iss,
	}
}
