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

type MultipleSentencesCreationRequest struct {
	Sentences []SentenceCreationRequest `json:"sentences"`
}

type SentenceCreation struct {
	Sentence    string
	LoginUserId uint64
}

type SentenceUpdateRequest struct {
	Id       uint64 `json:"id"`
	Sentence string `json:"sentence"`
}

type SentenceUpdate struct {
	Id          uint64
	Sentence    string
	LoginUserId uint64
}

type SentenceWithLink struct {
	Id               uint64
	Sentence         string
	SentenceWithLink string
	UserId           uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SentenceWithLinkResponse struct {
	Id               uint64 `json:"id"`
	Sentence         string `json:"sentence"`
	SentenceWithLink string `json:"sentence_with_link"`
	UserId           uint64 `json:"user_id"`
}

type SentencesCountResponse struct {
	Count uint64 `json:"count"`
}