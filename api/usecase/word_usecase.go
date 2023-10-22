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

func (wu *WordUsecase) GetAllWords(userId uint64) ([]model.WordResponse, error) {
	var wordResponses []model.WordResponse

	// TODO: userIdがログイン中のものと一致することを確認
	
	words, err := wu.wr.GetAllWords(userId)
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

func (wu *WordUsecase) GetWordById(userId uint64, wordId uint64) (model.WordResponse, error) {	
	// TODO: userIdがログイン中のものと一致することを確認

	word, err := wu.wr.GetWordById(wordId)
	if err != nil {
		return model.WordResponse{}, err
	}

	wordResponse := model.WordResponse{
		Id: word.Id,
		Word: word.Word,
		Memo: word.Memo,
		UserId: word.UserId,
	}

	return wordResponse, nil
}

func (wu *WordUsecase) CreateWord(userId uint64, req model.WordCreationRequest) (model.WordResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	newWord := model.WordCreation{
		Word: req.Word,
		Memo: req.Memo,
		UserId: userId,
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

func (wu *WordUsecase) DeleteWord(userId uint64, wordId uint64) (model.WordResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	deletedWord, err := wu.wr.DeleteWordById(wordId)
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

func (wu *WordUsecase) UpdateWord(userId uint64, req model.WordUpdateRequest) (model.WordResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	wordUpdate := model.WordUpdate{
		Id: req.Id,
		Word: req.Word,
		Memo: req.Memo,
		UserId: userId,
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
