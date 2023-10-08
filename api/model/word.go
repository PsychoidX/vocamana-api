package model

import "time"

type Word struct {
	Id        uint64
	Word      string
	Memo      string
	UserId    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WordResponse struct {
	Id     uint64   `json:"id"`
	Word   string `json:"word"`
	Memo   string `json:"memo"`
	UserId uint64   `json:"user_id"`
}

type WordCreationRequest struct {
	Word   string `json:"word"`
	Memo   string `json:"memo"`
}

type WordCreation struct {
	Word   string `json:"word"`
	Memo   string `json:"memo"`
	UserId uint64   `json:"user_id"`
}

type WordDeleteRequest struct {
	Id uint64 `json:"id"`
}
