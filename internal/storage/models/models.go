package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RoleModel struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FileModel struct {
	Id        int64     `json:"id"`
	Filename  string    `json:"filename"`
	Filepath  string    `json:"filepath"`
	SizeBytes int       `json:"size_bytes"`
	Creator   int64     `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

type TagModel struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ProfileModel struct {
	Id        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Bio       string    `json:"bio"`
	Avatar    FileModel `json:"avatar"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	Id         int64        `json:"id"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	Password   Password     `json:"-"`
	IsVerified bool         `json:"is_verifyed"`
	IsPrivate  bool         `json:"is_private"`
	Role       RoleModel    `json:"role"`
	Profile    ProfileModel `json:"profile"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type Password struct {
	Password string
	Hash     []byte
}

func (p *Password) Set(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash
	p.Password = pass

	return nil
}

func (p *Password) Compare(pass string) error {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(pass))
	return err
}

type PostModel struct {
	Id          int64       `json:"id"`
	Description string      `json:"description"`
	Creator     UserModel   `json:"creator"`
	Tags        []TagModel  `json:"tags"`
	Files       []FileModel `json:"files"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type PostFileModel struct {
	Id   int64     `json:"id"`
	Post PostModel `json:"post"`
	File FileModel `json:"file"`
}

type PostTagModel struct {
	Id   int64     `json:"id"`
	Post PostModel `json:"post"`
	Tag  TagModel  `json:"tag"`
}

type BanModel struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	WhyBanned string    `json:"why_banned"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentModel struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	CreatorID int64          `json:"creator_id"`
	PostID    int64          `json:"post_id"`
	ParentID  *int64         `json:"parent_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Post      *PostModel     `json:"post,omitempty"`
	User      *UserModel     `json:"user,omitempty"`
	Children  []CommentModel `json:"children,omitempty"`
}
