package model

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

type SentenceCreationRequest struct {
	Sentence string `json:"sentence"`
}

type SentenceCreation struct {
	Sentence string
	UserId   uint64
}

type SentenceUpdateRequest struct {
	Id       uint64 `json:"id"` 
	Sentence string `json:"sentence"`
}

type SentenceUpdate struct {
	Id       uint64
	Sentence string
	UserId   uint64
}