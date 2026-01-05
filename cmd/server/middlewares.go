package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
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

		ctx = context.WithValue(ctx, global_varables.USER_CTX, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) checkUserRoleMiddleware(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := helpers.GetUserFromContext(r)

		if user.Role.Id == 2 || user.Role.Name == global_varables.ADMIN_ROLE {
			next.ServeHTTP(w, r)
			return
		}

		s.forbiddenResponse(w, r, fmt.Errorf("only %s can access this route", role))
	}
}

func (s *server) checkAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := helpers.GetUserFromContext(r)

		userId, useridErr := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
		commentid, commentidErr := strconv.ParseInt(chi.URLParam(r, "commentid"), 10, 64)
		postid, postidErr := strconv.ParseInt(chi.URLParam(r, "postid"), 10, 64)

		hasAccess := false

		switch true {
		// check access to page
		case useridErr == nil:
			ctx := r.Context()
			targetUser, err := s.postgreStorage.UserStore.GetById(ctx, userId)
			if err != nil {
				switch {
				case errors.Is(err, global_varables.NOT_FOUND_ROW):
					s.notFoundResponse(w, r, err)
					return
				default:
					s.internalServerErrorResponse(w, r, err)
					return
				}
			}

			// checking is user page not private AND
			// is the one who wants to access it admin AND
			// is the one who wants to access same the same
			if !targetUser.IsPrivate &&
				user.Role.Name == global_varables.ADMIN_ROLE &&
				user.Id == targetUser.Id {
				hasAccess = true
			}

		// check access to comment
		case commentidErr == nil:
			ctx := r.Context()
			comment, err := s.postgreStorage.CommentStore.GetById(ctx, commentid)

			if err != nil {
				switch {
				case errors.Is(err, global_varables.NOT_FOUND_ROW):
					s.notFoundResponse(w, r, err)
					return
				default:
					s.internalServerErrorResponse(w, r, err)
					return
				}
			}

			// checking is user page not private AND
			// is the one who wants to access it admin AND
			// is the one who wants to access same the same
			if !comment.Creator.IsPrivate &&
				user.Role.Name == global_varables.ADMIN_ROLE &&
				user.Id == comment.Creator.Id {
				hasAccess = true
			}

		// check access to post
		case postidErr == nil:
			ctx := r.Context()
			post, err := s.postgreStorage.PostStore.GetById(ctx, postid)

			if err != nil {
				switch {
				case errors.Is(err, global_varables.NOT_FOUND_ROW):
					s.notFoundResponse(w, r, err)
					return
				default:
					s.internalServerErrorResponse(w, r, err)
					return
				}
			}

			// checking is user page not private AND
			// is the one who wants to access it admin AND
			// is the one who wants to access same the same
			if !post.Creator.IsPrivate &&
				user.Role.Name == global_varables.ADMIN_ROLE &&
				user.Id == post.Creator.Id {
				hasAccess = true
			}
		default:
			hasAccess = false
		}

		if !hasAccess {
			s.forbiddenResponse(w, r, fmt.Errorf("this user is private and you cannot have access on it"))
			return
		}

		next.ServeHTTP(w, r)

	}
}

func (s *server) checkIsUserVerifiedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.GetUserFromContext(r)

		if !user.IsVerified {
			s.forbiddenResponse(w, r, fmt.Errorf("you are not verified"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
