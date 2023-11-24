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

func (wu *WordUsecase) GetAllNotations(loginUserId, wordId uint64) ([]model.Notation, error) {
	// wordIdの所有者がloginUserIdの場合ゼロ値を返す
	isWordOwner, err := wu.wr.IsWordOwner(wordId, loginUserId)
	if err != nil {
		return []model.Notation{}, err
	}
	if !isWordOwner {
		return []model.Notation{}, nil
	}

	notations, err := wu.nr.GetAllNotations(wordId)
	if err != nil {
		return []model.Notation{}, err
	}

	return notations, nil
}

func (wu *WordUsecase) CreateNotation(notationCreation model.NotationCreation) (model.Notation, error) {
	loginUserId := notationCreation.LoginUserId

	// 追加先のWordIdの所有者がloginUserIdでない場合何もしない
	isWordOwner, err := wu.wr.IsWordOwner(notationCreation.WordId, loginUserId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}

	createdNotation, err := wu.nr.InsertNotation(notationCreation)
	if err != nil {
		return model.Notation{}, err
	}

	// 既存のSentenceに追加されたWord含まれればsentences_wordsに追加
	AssociateWordWithAllSentences(
		loginUserId,
		createdNotation.WordId,
		wu.wr,
		wu.sr,
		wu.swr,
		wu.nr,
	)

	return createdNotation, nil
}

func (wu *WordUsecase) UpdateNotation(notationUpdate model.NotationUpdate) (model.Notation, error) {
	notation, err := wu.nr.GetNotationById(notationUpdate.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 更新対象のNotationが存在しない場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	// WordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := wu.wr.IsWordOwner(notation.WordId, notationUpdate.LoginUserId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}

	updatedNotation, err := wu.nr.UpdateNotation(notationUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	ReAssociateWordWithAllSentences(notationUpdate.LoginUserId, notation.WordId, wu.wr, wu.sr, wu.swr, wu.nr)

	return updatedNotation, nil
}

func (wu *WordUsecase) DeleteNotation(loginUserId, notationId uint64) (model.Notation, error) {
	notation, err := wu.nr.GetNotationById(notationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// 削除対象のNotationが存在しない場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	// WordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := wu.wr.IsWordOwner(notation.WordId, loginUserId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}

	deletedNotation, err := wu.nr.DeleteNotationById(notationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが削除されなかった場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	ReAssociateWordWithAllSentences(loginUserId, notation.WordId, wu.wr, wu.sr, wu.swr, wu.nr)

	return deletedNotation, nil
}