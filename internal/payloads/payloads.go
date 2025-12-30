package payloads

import (
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type CreateUserPayload struct {
	Username string            `json:"username" validate:"required,min=3,max=255,alphanum"`
	Fullname string            `json:"fullname" validate:"required,min=8,max=255"`
	Email    string            `json:"email" validate:"required,email,max=255"`
	Password string            `json:"password" validate:"required,min=8,max=255"`
	Bio      string            `json:"bio" validate:"max=512"`
	Avatar   CreateFilePayload `json:"avatar"`
}

type CreateFilePayload struct {
	Filename  string `json:"filename"`
	Filepath  string `json:"filepath"`
	SizeBytes int    `json:"size_bytes"`
	Creator   int64  `json:"creator"`
}

type CreateBanPayload struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	WhyBanned string `json:"why_banned" validate:"required, min=8"`
}

type CreatePostPayload struct {
	Description string              `json:"description" validate:"max=1024"`
	Creator     int64               `json:"creator" validate:"required,numeric,min=1"`
	Tags        []string            `json:"tags" validate="max=30,dive,required,min=1,max=255"`
	Files       []CreateFilePayload `json:"files"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserWithToken struct {
	User  models.UserModel `json:"user"`
	Token string           `json:"token"`
}
