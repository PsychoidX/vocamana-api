package usecase

import (
	"api/model"
	"api/repository"
)

type WordUsecase struct {
	wr repository.IWordRepository
}

func NewWordUsecase(wr repository.IWordRepository) *WordUsecase {
	return &WordUsecase{wr}
}

func (wu *WordUsecase) GetAllWords() ([]model.WordResponse, error) {
	// TODO
	return []model.WordResponse{}, nil
}

func (wu *WordUsecase) GetWordById(id uint) (model.WordResponse, error) {
	// TODO
	return model.WordResponse{}, nil
}

func (wu *WordUsecase) CreateWord(newWord model.WordRegistration) (model.WordResponse, error) {
	createdWord, err := wu.wr.InsertWord(newWord)
	if err != nil {
		return model.WordResponse{}, err
	}

	createdWordResponse := model.WordResponse{
		Id:     createdWord.Id,
		Word:   createdWord.Word,
		Memo:   createdWord.Memo,
		UserId: createdWord.UserId,
	}

	return createdWordResponse, nil
}

