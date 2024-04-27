package model

type WordIdsRequest struct {
	WordIds []uint64 `json:"word_ids"`
}

type WordIds struct {
	WordIds []uint64
}

type WordIdsResponse struct {
	WordIds []uint64 `json:"word_ids"`
}