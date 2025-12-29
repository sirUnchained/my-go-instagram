package helpers

import (
	"net/http"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

func GetUserFromContext(r *http.Request) *models.UserModel {
	user, _ := r.Context().Value(global_varables.USER_CTX).(models.UserModel)
	return &user
}

func GetUserByIdFromContext(r *http.Request) *models.UserModel {
	user, _ := r.Context().Value(global_varables.TARGET_USER_CTX).(models.UserModel)
	return &user
}
