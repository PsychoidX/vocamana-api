package model

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

// Word作成時に必要な入力項目
type WordRegistrationInput struct {
	Word   string `json:"word"`
	Memo   string `json:"memo"`
}

// Word作成時に必要な項目
type WordRegistration struct {
	Word   string `json:"word"`
	Memo   string `json:"memo"`
	UserId uint   `json:"user_id"`
}