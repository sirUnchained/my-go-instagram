package models

import "time"

type ProfileModel struct {
	Id        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	UpdatedAt time.Time `json:"updated_at"`
}
