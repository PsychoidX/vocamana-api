package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
)

type WordUsecase struct {
	wr  repository.IWordRepository
	sr  repository.ISentenceRepository
	swr repository.ISentencesWordsRepository
	nr  repository.INotationRepository
}

func NewWordUsecase(
	wr repository.IWordRepository,
	sr repository.ISentenceRepository,
	swr repository.ISentencesWordsRepository,
	nr repository.INotationRepository,
) *WordUsecase {
	return &WordUsecase{wr, sr, swr, nr}
}

func (wu *WordUsecase) GetAllWords(loginUserId uint64) ([]model.Word, error) {
	words, err := wu.wr.GetAllWords(loginUserId)
	if err != nil {
		return []model.Word{}, err
	}

	return words, nil
}

func (wu *WordUsecase) GetWordById(loginUserId, wordId uint64) (model.Word, error) {
	word, err := wu.wr.GetWordById(loginUserId, wordId)
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
	loginUserId := wordCreation.LoginUserId

	createdWord, err := wu.wr.InsertWord(wordCreation)
	if err != nil {
		return model.Word{}, err
	}

	// 既存のSentenceに追加されたWord含まれればsentences_wordsに追加
	AssociateWordWithAllSentences(
		loginUserId,
		createdWord.Id,
		wu.wr,
		wu.sr,
		wu.swr,
		wu.nr,
	)

	return createdWord, nil
}

func (wu *WordUsecase) CreateMultipleWords(wordCreations []model.WordCreation) ([]model.Word, error) {
	// TODO 1件でも失敗したらロールバックする実装に変更
	var createdWords []model.Word
	for _, wordCreation := range wordCreations {
		createdWord, err := wu.CreateWord(wordCreation)
		if err != nil {
			return []model.Word{}, err
		}

		createdWords = append(createdWords, createdWord)
	}

	return createdWords, nil
}

func (wu *WordUsecase) DeleteWord(loginUserId, wordId uint64) (model.Word, error) {
	deletedWord, err := wu.wr.DeleteWordById(loginUserId, wordId)
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

	ReAssociateWordWithAllSentences(wordUpdate.LoginUserId, wordUpdate.Id, wu.wr, wu.sr, wu.swr, wu.nr)

	return updatedWord, nil
}