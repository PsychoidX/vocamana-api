package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
)

type WordUsecase struct {
	wr repository.IWordRepository
}

func NewWordUsecase(wr repository.IWordRepository) *WordUsecase {
	return &WordUsecase{wr}
}

func (wu *WordUsecase) GetAllWords(userId uint64) ([]model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	words, err := wu.wr.GetAllWords(userId)
	if err != nil {
		return []model.Word{}, err
	}

	return words, nil
}

func (wu *WordUsecase) GetWordById(userId uint64, wordId uint64) (model.Word, error) {	
	// TODO: userIdがログイン中のものと一致することを確認

	word, err := wu.wr.GetWordById(userId, wordId)
	if err != nil {
		if err == sql.ErrNoRows {
			// マッチするレコードが無い場合
			// Wordのゼロ値を返す
			return model.Word{}, nil
		}
		return model.Word{}, err
	}

	return word, nil
}

func (wu *WordUsecase) CreateWord(wordCreation model.WordCreation) (model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	createdWord, err := wu.wr.InsertWord(wordCreation)
	if err != nil {
		return model.Word{}, err
	}

	return createdWord, nil
}

func (wu *WordUsecase) DeleteWord(userId uint64, wordId uint64) (model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	deletedWord, err := wu.wr.DeleteWordById(userId, wordId)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが削除されなかった場合
			// Wordのゼロ値を返す
			return model.Word{}, nil
		}
		
		return model.Word{}, err
	}

	return deletedWord, nil
}

func (wu *WordUsecase) UpdateWord(wordUpdate model.WordUpdate) (model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	updatedWord, err := wu.wr.UpdateWord(wordUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Wordのゼロ値を返す
			return model.Word{}, nil
		}

		return model.Word{}, err
	}

	return updatedWord, nil
}
