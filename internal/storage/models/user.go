package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Password   Password  `json:"-"`
	IsVerified bool      `json:"is_verifyed"`
	Role       RoleModel `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UreatedAt  time.Time `json:"updated_at"`
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
