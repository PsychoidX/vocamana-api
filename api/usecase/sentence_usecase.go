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

func (su *SentenceUsecase) GetAllSentences(userId uint64) ([]model.SentenceResponse, error) {
	var sentenceResponses []model.SentenceResponse

	// TODO: userIdがログイン中のものと一致することを確認
	
	sentences, err := su.sr.GetAllSentences(userId)
	if err != nil {
		return []model.SentenceResponse{}, err
	}

	for _, sentence := range sentences {
		sentenceResponse := model.SentenceResponse{
			Id: sentence.Id,
			Sentence: sentence.Sentence,
			UserId: sentence.UserId,
		}
		sentenceResponses = append(sentenceResponses, sentenceResponse)
	}

	return sentenceResponses, nil
}

func (su *SentenceUsecase) GetSentenceById(userId uint64, sentenceId uint64) (model.SentenceResponse, error) {	
	// TODO: userIdがログイン中のものと一致することを確認

	sentence, err := su.sr.GetSentenceById(userId, sentenceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// マッチするレコードが無い場合
			// SentenceResponseのゼロ値を返す
			return model.SentenceResponse{}, nil
		}
		
		return model.SentenceResponse{}, err
	}

	sentenceResponse := model.SentenceResponse{
		Id: sentence.Id,
		Sentence: sentence.Sentence,
		UserId: sentence.UserId,
	}

	return sentenceResponse, nil
}

func (su *SentenceUsecase) CreateSentence(userId uint64, req model.SentenceCreationRequest) (model.SentenceResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	newSentence := model.SentenceCreation{
		Sentence: req.Sentence,
		UserId: userId,
	}

	createdSentence, err := su.sr.InsertSentence(newSentence)
	if err != nil {
		return model.SentenceResponse{}, err
	}

	createdSentenceResponse := model.SentenceResponse{
		Id:     createdSentence.Id,
		Sentence:   createdSentence.Sentence,
		UserId: createdSentence.UserId,
	}

	return createdSentenceResponse, nil
}

func (su *SentenceUsecase) UpdateSentence(userId uint64, req model.SentenceUpdateRequest) (model.SentenceResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	sentenceUpdate := model.SentenceUpdate{
		Id: req.Id,
		Sentence: req.Sentence,
		UserId: userId,
	}

	updatedSentence, err := su.sr.UpdateSentence(sentenceUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// SentenceResponseのゼロ値を返す
			return model.SentenceResponse{}, nil
		}

		return model.SentenceResponse{}, err
	}

	updatedSentenceResponse := model.SentenceResponse{
		Id:     updatedSentence.Id,
		Sentence:   updatedSentence.Sentence,
		UserId: updatedSentence.UserId,
	}

	return updatedSentenceResponse, nil
}

func (su *SentenceUsecase) DeleteSentence(userId uint64, sentenceId uint64) (model.SentenceResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	deletedSentence, err := su.sr.DeleteSentenceById(userId, sentenceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが削除されなかった場合
			// SentenceResponseのゼロ値を返す
			return model.SentenceResponse{}, nil
		}

		return model.SentenceResponse{}, err
	}

	deletedSentenceResponse := model.SentenceResponse{
		Id:     deletedSentence.Id,
		Sentence:   deletedSentence.Sentence,
		UserId: deletedSentence.UserId,
	}

	return deletedSentenceResponse, nil
}

func (su *SentenceUsecase) AssociateSentenceWithWords(userId uint64, sentenceId uint64, wordIds model.WordIdsRequest) (model.WordIdsResponse, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	// sentenceeIdの所有者がuserIdでない場合何もしない
	isSentenceOwner, err := su.sr.IsSentenceOwner(sentenceId, userId)
	if err != nil {
		return model.WordIdsResponse{}, err
	}
	if !isSentenceOwner {
		return model.WordIdsResponse{}, nil
	}

	var associatedWordIds []uint64

	for _, wordId := range wordIds.WordIds {
		// wordIdがuserIdのものでない場合continue
		isWordOwner, err := su.wr.IsWordOwner(wordId, userId)
		if err != nil {
			return model.WordIdsResponse{}, err
		}
		if !isWordOwner {
			continue
		}

		err = su.swr.AssociateSentenceWithWord(sentenceId, wordId)
		if err != nil {
			return model.WordIdsResponse{}, err
		}
		associatedWordIds = append(associatedWordIds, wordId)
	}

	wordIdsResponse := model.WordIdsResponse{
		WordIds: associatedWordIds,
	}

	return wordIdsResponse, nil
}