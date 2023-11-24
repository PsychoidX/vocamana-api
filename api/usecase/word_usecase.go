package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
	"strings"
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
	wu.AssociateWordWithAllSentences(loginUserId, createdWord.Id)

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

	wu.ReAssociateWordWithAllSentences(wordUpdate.LoginUserId, wordUpdate.Id)

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
	wu.AssociateWordWithAllSentences(loginUserId, createdNotation.WordId)

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

	wu.ReAssociateWordWithAllSentences(notationUpdate.LoginUserId, notation.WordId)

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

	wu.ReAssociateWordWithAllSentences(loginUserId, notation.WordId)

	return deletedNotation, nil
}

func (wu *WordUsecase)AssociateWordWithAllSentences(userId, wordId uint64) ([]model.Sentence, error) {
	// userIdに紐づく全Sentenceに対し、
	// Sentenceの中にwordIdのWordまたはNotationが含まれればsentences_wordsにレコード追加

	// TODO: userIdがログイン中のものと一致することを確認

	// wordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := wu.wr.IsWordOwner(wordId, userId)
	if err != nil {
		return []model.Sentence{}, err
	}
	if !isWordOwner {
		return []model.Sentence{}, nil
	}

	word, err := wu.GetWordById(userId, wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	userSentences, err := wu.sr.GetAllSentences(userId)
	if err != nil {
		return []model.Sentence{}, err
	}

	var associatedSentences []model.Sentence
	for _, sentence := range userSentences {
		// Sentence中にWordが含まれるか判定
		if strings.Contains(sentence.Sentence, word.Word) {
			err = wu.swr.AssociateSentenceWithWord(sentence.Id, word.Id)
			if err != nil {
				return []model.Sentence{}, err
			}
			associatedSentences = append(associatedSentences, sentence)
			// Sentenceの中にWordが含まれる場合、
			// continueし、Notationが含まれるかの判定はしない
			continue
		}

		// Sentence中にNotationが含まれるか判定
		notations, err := wu.nr.GetAllNotations(word.Id)
		if err != nil {
			return []model.Sentence{}, err
		}

		for _, notation := range notations {
			if strings.Contains(sentence.Sentence, notation.Notation) {
				err = wu.swr.AssociateSentenceWithWord(sentence.Id, word.Id)
				if err != nil {
					return []model.Sentence{}, err
				}
				associatedSentences = append(associatedSentences, sentence)
				// Sentenceの中にNotationが含まれる場合、
				// sentences_wordsに2つ目のレコードが追加されないようbreak
				break
			}
		}
	}

	return associatedSentences, nil
}

func (wu *WordUsecase)ReAssociateWordWithAllSentences(loginUserId, wordId uint64) error {
	// wordIdで指定されるWordと、全Sentenceのsentences_wordsを再構築
	// sentences_wordsからwordIdのレコードを全削除し、もう一度追加しなおす

	// sentenceIdの所有者がloginUserIdでない場合何もしない
	isWordOwner, err := wu.wr.IsWordOwner(wordId, loginUserId)
	if err != nil {
		return err
	}
	if !isWordOwner {
		return nil
	}

	// TODO: 削除～再追加はトランザクション内で行う

	// sentences_wordsからwordIdのレコードを全削除
	err = wu.swr.DeleteAllAssociationByWordId(wordId)

	// sentences_wordsに再追加
	wu.AssociateWordWithAllSentences(loginUserId, wordId)

	return nil
}