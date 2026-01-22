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

type UnbanPayload struct {
	Email string `json:"email" validate:"required,email"`
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

type CreateCommentPayload struct {
	Content   string `json:"content" validate:"required,min=5,max=2048"`
	CreatorID int64  `json:"creator_id" validate:"required,numeric,min=1"`
	PostID    int64  `json:"post_id" validate:"required,numeric,min=1"`
	ParentID  *int64 `json:"parent_id" validate:"numeric,min=1"`
}

type CreateReportPayload struct {
	CreatorID int    `json:"creator_id" validate:"required,numeric,min=1"`
	PostID    int    `json:"post_id" validate:"omitempty,numeric,min=1"`
	CommentID int    `json:"comment_id" validate:"omitempty,numeric,min=1"`
	Reason    string `json:"reason" validate:"required,oneof=spam_report porn_content racist_content other"`
	Content   string `json:"content" validate:"max=500"`
}
