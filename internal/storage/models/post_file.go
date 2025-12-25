package models

type PostFileModel struct {
	Id   int64     `json:"id"`
	Post PostModel `json:"post"`
	File FileModel `json:"file"`
}
