package model

import "time"

type Notation struct {
	Id        uint64
	WordId    uint64
	Notation  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotationResponse struct {
	Id       uint64 `json:"id"`
	WordId   uint64 `json:"word_id"`
	Notation string `json:"notation"`
}

type NotationCreationRequest struct {
	Notation string `json:"notation"`
}

type NotationCreation struct {
	WordId       uint64
	Notation     string
	LoginUserId  uint64
}

type NotationUpdateRequest struct {
	Notation string `json:"notation"`
}

type NotationUpdate struct {
	Id          uint64
	WordId      uint64
	Notation    string
	LoginUserId uint64
}