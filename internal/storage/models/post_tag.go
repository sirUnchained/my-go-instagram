package models

type PostTagModel struct {
	Id   int64     `json:"id"`
	Post PostModel `json:"post"`
	Tag  TagModel  `json:"tag"`
}
