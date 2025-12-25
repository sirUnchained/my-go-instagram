package models

import "time"

type FileModel struct {
	Id        int64     `json:"id"`
	Filename  string    `json:"filename"`
	Filepath  string    `json:"filepath"`
	SizeBytes int       `json:"size_bytes"`
	Creator   UserModel `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}
