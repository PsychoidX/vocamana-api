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

func (wu *WordUsecase) GetAllWords(user_id uint64) ([]model.WordResponse, error) {
	var wordResponses []model.WordResponse
	
	words, err := wu.wr.GetAllWords(user_id)
	if err != nil {
		return []model.WordResponse{}, err
	}

	for _, word := range words {
		wordResponse := model.WordResponse{
			Id: word.Id,
			Word: word.Word,
			Memo: word.Memo,
			UserId: word.UserId,
		}
		wordResponses = append(wordResponses, wordResponse)
	}

	return wordResponses, nil
}

func (wu *WordUsecase) GetWordById(id uint64) (model.WordResponse, error) {
	// TODO
	return model.WordResponse{}, nil
}

func (wu *WordUsecase) CreateWord(req model.WordCreationRequest) (model.WordResponse, error) {
	newWord := model.WordCreation{
		Word: req.Word,
		Memo: req.Memo,
		UserId: 1, // TODO リクエストに含める
	}

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

func (wu *WordUsecase) DeleteWord(id uint64) (model.WordResponse, error) {
	deletedWord, err := wu.wr.DeleteWordById(id)
	if err != nil {
		return model.WordResponse{}, err
	}

	deletedWordResponse := model.WordResponse{
		Id:     deletedWord.Id,
		Word:   deletedWord.Word,
		Memo:   deletedWord.Memo,
		UserId: deletedWord.UserId,
	}

	return deletedWordResponse, nil
}

func (wu *WordUsecase) UpdateWord(req model.WordUpdateRequest) (model.WordResponse, error) {
	wordUpdate := model.WordUpdate{
		Id: req.Id,
		Word: req.Word,
		Memo: req.Memo,
		UserId: 1, // TODO リクエストに含める
	}

	updatedWord, err := wu.wr.UpdateWord(wordUpdate)
	if err != nil {
		return model.WordResponse{}, err
	}

	updatedWordResponse := model.WordResponse{
		Id:     updatedWord.Id,
		Word:   updatedWord.Word,
		Memo:   updatedWord.Memo,
		UserId: updatedWord.UserId,
	}

	return updatedWordResponse, nil
}
