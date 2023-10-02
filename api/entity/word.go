package entity

import "time"

type Word struct {
	Id        uint
	Word      string
	Memo      string
	UserId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WordResponse struct {
	Id     uint   `json:"id"`
	Word   string `json:"word"`
	Memo   string `json:"memo"`
	UserId uint   `json:"user_id"`
}