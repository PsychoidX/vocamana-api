package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
	"strings"
)

type SentenceUsecase struct {
	sr  repository.ISentenceRepository
	wr  repository.IWordRepository
	swr repository.ISentencesWordsRepository
	nr  repository.INotationRepository
}

func NewSentenceUsecase(
	sr repository.ISentenceRepository,
	wr repository.IWordRepository,
	swr repository.ISentencesWordsRepository,
	nr repository.INotationRepository,
) *SentenceUsecase {
	return &SentenceUsecase{sr, wr, swr, nr}
}

func (su *SentenceUsecase) GetAllSentences(loginUserId uint64) ([]model.Sentence, error) {
	sentences, err := su.sr.GetAllSentences(loginUserId)
	if err != nil {
		return []model.Sentence{}, err
	}

	return sentences, nil
}

func (su *SentenceUsecase) GetSentenceById(loginUserId uint64, sentenceId uint64) (model.Sentence, error) {
	sentence, err := su.sr.GetSentenceById(loginUserId, sentenceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// マッチするレコードが無い場合
			// Sentenceのゼロ値を返す
			return model.Sentence{}, nil
		}

		return model.Sentence{}, err
	}

	return sentence, nil
}

func (su *SentenceUsecase) CreateSentence(sentenceCreation model.SentenceCreation) (model.Sentence, error) {
	loginUserId := sentenceCreation.LoginUserId

	createdSentence, err := su.sr.InsertSentence(sentenceCreation)
	if err != nil {
		return model.Sentence{}, err
	}

	// 追加されたSentenceに既存のWordが含まれればsentences_wordsに追加
	su.AssociateSentenceWithAllWords(loginUserId, createdSentence.Id)

	return createdSentence, nil
}

func (su *SentenceUsecase) CreateMultipleSentences(sentenceCreations []model.SentenceCreation) ([]model.Sentence, error) {
	// TODO 1件でも失敗したらロールバックする実装に変更
	var createdSentences []model.Sentence
	for _, sentenceCreation := range sentenceCreations {
		createdSentence, err := su.CreateSentence(sentenceCreation)
		if err != nil {
			return []model.Sentence{}, err
		}

		createdSentences = append(createdSentences, createdSentence)
	}
	
	return createdSentences, nil
}

func (su *SentenceUsecase) UpdateSentence(sentenceUpdate model.SentenceUpdate) (model.Sentence, error) {
	updatedSentence, err := su.sr.UpdateSentence(sentenceUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Sentenceのゼロ値を返す
			return model.Sentence{}, nil
		}

		return model.Sentence{}, err
	}

	su.ReAssociateSentenceWithAllWords(sentenceUpdate.LoginUserId, sentenceUpdate.Id)

	return updatedSentence, nil
}

func (su *SentenceUsecase) DeleteSentence(loginUserId uint64, sentenceId uint64) (model.Sentence, error) {
	deletedSentence, err := su.sr.DeleteSentenceById(loginUserId, sentenceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが削除されなかった場合
			// Sentenceのゼロ値を返す
			return model.Sentence{}, nil
		}

		return model.Sentence{}, err
	}

	return deletedSentence, nil
}

func (su *SentenceUsecase) GetAssociatedWordsBySentenceId(loginUserId uint64, sentenceId uint64) ([]model.Word, error) {
	// sentenceIdの所有者がloginUserIdでない場合ゼロ値を返す
	isSentenceOwner, err := su.sr.IsSentenceOwner(sentenceId, loginUserId)
	if err != nil {
		return []model.Word{}, err
	}
	if !isSentenceOwner {
		return []model.Word{}, nil
	}

	words, err := su.swr.GetUserAssociatedWordsBySentenceId(sentenceId)
	if err != nil {
		return []model.Word{}, err
	}

	// リポジトリの返り値のuserIdを検証
	userWords := []model.Word{}
	for _, word := range words {
		// wordの所有者がloginUserIdでない場合continue
		isWordOwner, err := su.wr.IsWordOwner(word.Id, loginUserId)
		if err != nil {
			return []model.Word{}, err
		}
		if !isWordOwner {
			continue
		}

		userWords = append(userWords, word)
	}

	return userWords, nil
}

func (su *SentenceUsecase) AssociateSentenceWithAllWords(loginUserId , sentenceId uint64) ([]model.Word, error) {
	// loginUserIdに紐づく全Wordに対し、
	// sentenceIdのSentence中にWordまたはNotationが含まれればsentences_wordsにレコード追加

	// sentenceIdの所有者がloginUserIdでない場合何もしない
	isSentenceOwner, err := su.sr.IsSentenceOwner(sentenceId, loginUserId)
	if err != nil {
		return []model.Word{}, err
	}
	if !isSentenceOwner {
		return []model.Word{}, nil
	}

	sentence, err := su.GetSentenceById(loginUserId, sentenceId)
	if err != nil {
		return []model.Word{}, err
	}

	userWords, err := su.wr.GetAllWords(loginUserId)
	if err != nil {
		return []model.Word{}, err
	}

	var associatedWords []model.Word
	for _, word := range userWords {
		// Sentence中にWordが含まれる場合
		// sentences_wordsにレコード追加
		if strings.Contains(sentence.Sentence, word.Word) {
			err = su.swr.AssociateSentenceWithWord(sentence.Id, word.Id)
			if err != nil {
				return []model.Word{}, err
			}
			associatedWords = append(associatedWords, word)
			// Sentenceの中にWordが含まれる場合、
			// continueし、Notationが含まれるかの判定はしない
			continue
		}

		// Sentence中にNotationが含まれる場合
		// sentences_wordsにレコード追加
		notations, err := su.nr.GetAllNotations(word.Id)
		if err != nil {
			return []model.Word{}, err
		}

		for _, notation := range notations {
			if strings.Contains(sentence.Sentence, notation.Notation) {
				err = su.swr.AssociateSentenceWithWord(sentence.Id, word.Id)
				if err != nil {
					return []model.Word{}, err
				}
				associatedWords = append(associatedWords, word)
				// Sentenceの中にNotationが含まれる場合、
				// sentences_wordsに2つ目のレコードが追加されないようbreak
				break
			}
		}
	}

	return associatedWords, nil
}

func (su *SentenceUsecase) ReAssociateSentenceWithAllWords(loginUserId, sentenceId uint64) error {
	// sentenceIdと全Wordのsentences_wordsを再構築
	// sentences_wordsからsentenceIdのレコードを全削除し、もう一度追加しなおす

	// sentenceIdの所有者がloginUserIdでない場合何もしない
	isSentenceOwner, err := su.sr.IsSentenceOwner(sentenceId, loginUserId)
	if err != nil {
		return err
	}
	if !isSentenceOwner {
		return nil
	}

	// TODO: 削除～再追加はトランザクション内で行う

	// sentences_wordsからsentenceIdのレコードを全削除
	err = su.swr.DeleteAllAssociationBySentenceId(sentenceId)

	// sentences_wordsに再追加
	su.AssociateSentenceWithAllWords(loginUserId, sentenceId)

	return nil
}

func (su *SentenceUsecase) GetSentencesCount(loginUserId uint64) (uint64, error) {
	return su.sr.GetSentencesCount(loginUserId)
}