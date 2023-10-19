package model

type WordIdsRequest struct {
	WordIds []uint64 `json:"word_ids"`
}

type WordIdsResponse struct {
	WordIds []uint64 `json:"word_ids"`
}