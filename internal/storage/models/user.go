package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
