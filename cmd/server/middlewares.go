package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
)

func (s *server) checkUserTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" || !strings.Contains(auth, "Bearer ") {
			s.unauthorizedResponse(w, r, fmt.Errorf("invalid token"))
			return
		}

		token := strings.Split(auth, " ")[1]
		jwtToken, err := s.auth.ValidateToken(token)
		if err != nil {
			s.unauthorizedResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userid_Str, err := claims.GetSubject()
		if err != nil {
			s.unauthorizedResponse(w, r, err)
			return
		}

		userid, err := strconv.ParseInt(userid_Str, 10, 64)
		ctx := r.Context()
		user, err := s.postgreStorage.UserStore.GetById(ctx, userid)
		if err != nil {
			switch {
			case errors.Is(err, global_varables.NOT_FOUND_ROW):
				s.unauthorizedResponse(w, r, fmt.Errorf("the id that token provided dose not exists in postgres storage."))
				return
			default:
				s.internalServerErrorResponse(w, r, err)
				return
			}
		}

		fmt.Printf("%++v\n\n", user)

		ctx = context.WithValue(ctx, global_varables.USER_CTX, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
