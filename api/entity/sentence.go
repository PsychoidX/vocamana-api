package entity

import "time"

type Sentence struct {
	Id        uint
	Sentence  string
	UserId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SentenceResponse struct {
	Id       uint   `json:"id"`
	Sentence string `json:"sentence"`
	UserId   uint   `json:"user_id"`
}