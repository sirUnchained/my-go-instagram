package models

import "time"

type PostModel struct {
	Id          int64       `json:"id"`
	Description string      `json:"description"`
	Creator     UserModel   `json:"creator"`
	Tags        []TagModel  `json:"tags"`
	Files       []FileModel `json:"files"`
	CreatedAt   time.Time   `json:"created_at"`
	UreatedAt   time.Time   `json:"updated_at"`
}
