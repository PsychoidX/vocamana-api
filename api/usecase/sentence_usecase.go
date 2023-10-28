package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
)

type SentenceUsecase struct {
	sr repository.ISentenceRepository
	wr repository.IWordRepository
	swr repository.ISentencesWordsRepository
}

func NewSentenceUsecase(
		sr repository.ISentenceRepository,
		wr repository.IWordRepository,
		swr repository.ISentencesWordsRepository,
	) *SentenceUsecase {
	return &SentenceUsecase{sr, wr, swr}
}

func (su *SentenceUsecase) GetAllSentences(userId uint64) ([]model.Sentence, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	sentences, err := su.sr.GetAllSentences(userId)
	if err != nil {
		return []model.Sentence{}, err
	}

	return sentences, nil
}

func (su *SentenceUsecase) GetSentenceById(userId uint64, sentenceId uint64) (model.Sentence, error) {	
	// TODO: userIdがログイン中のものと一致することを確認

	sentence, err := su.sr.GetSentenceById(userId, sentenceId)
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
	// TODO: userIdがログイン中のものと一致することを確認

	createdSentence, err := su.sr.InsertSentence(sentenceCreation)
	if err != nil {
		return model.Sentence{}, err
	}

	return createdSentence, nil
}

func (su *SentenceUsecase) UpdateSentence(sentenceUpdate model.SentenceUpdate) (model.Sentence, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	updatedSentence, err := su.sr.UpdateSentence(sentenceUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Sentenceのゼロ値を返す
			return model.Sentence{}, nil
		}

		return model.Sentence{}, err
	}

	return updatedSentence, nil
}

func (su *SentenceUsecase) DeleteSentence(userId uint64, sentenceId uint64) (model.Sentence, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	deletedSentence, err := su.sr.DeleteSentenceById(userId, sentenceId)
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

func (su *SentenceUsecase) AssociateSentenceWithWords(userId uint64, sentenceId uint64, wordIds model.WordIds) (model.WordIds, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	// sentenceeIdの所有者がuserIdでない場合何もしない
	isSentenceOwner, err := su.sr.IsSentenceOwner(sentenceId, userId)
	if err != nil {
		return model.WordIds{}, err
	}
	if !isSentenceOwner {
		return model.WordIds{}, nil
	}

	var associatedWordIds []uint64

	for _, wordId := range wordIds.WordIds {
		// wordIdの所有者がuserIdでない場合continue
		isWordOwner, err := su.wr.IsWordOwner(wordId, userId)
		if err != nil {
			return model.WordIds{}, err
		}
		if !isWordOwner {
			continue
		}

		err = su.swr.AssociateSentenceWithWord(sentenceId, wordId)
		if err != nil {
			return model.WordIds{}, err
		}
		associatedWordIds = append(associatedWordIds, wordId)
	}

	resultWordIds := model.WordIds{
		WordIds: associatedWordIds,
	}

	return resultWordIds, nil
}