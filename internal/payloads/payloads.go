package payloads

import "github.com/sirUnchained/my-go-instagram/internal/storage/models"

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,max=255"`
	Fullname string `json:"fullname" validate:"required,min=8,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserWithToken struct {
	User  models.UserModel `json:"user"`
	Token string           `json:"token"`
}
