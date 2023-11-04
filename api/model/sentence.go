package model

import "time"

type Sentence struct {
	Id        uint64
	Sentence  string
	UserId    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SentenceResponse struct {
	Id       uint64 `json:"id"`
	Sentence string `json:"sentence"`
	UserId   uint64 `json:"user_id"`
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

type SentenceWithLink struct {
	Id               uint64
	SentenceWithLink string
	UserId           uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
type SentenceWIthLinkResponse struct {
	Id               uint64 `json:"id"`
	SentenceWithLink string `json:"sentence"`
	UserId           uint64 `json:"user_id"`
}
