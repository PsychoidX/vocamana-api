package usecase

import (
	"api/model"
	"api/repository"
)

type IWordUsecase interface {
	GetAllWords() ([]model.WordResponse, error)
	GetWordById(id uint) (model.WordResponse, error)
}

type wordUsecase struct {
	wr repository.IWordRepository
}

func NewWordUsecase(wr repository.IWordRepository) IWordUsecase {
	return &wordUsecase{wr}
}

func (wu *wordUsecase) GetAllWords() ([]model.WordResponse, error) {
	// TODO
	return []model.WordResponse{}, nil
}

func (wu *wordUsecase) GetWordById(id uint) (model.WordResponse, error) {
	// TODO
	return model.WordResponse{}, nil
}