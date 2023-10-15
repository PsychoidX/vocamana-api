package usecase

import (
	"api/model"
	"api/repository"
)

type SentenceUsecase struct {
	sr repository.ISentenceRepository
}

func NewSentenceUsecase(sr repository.ISentenceRepository) *SentenceUsecase {
	return &SentenceUsecase{sr}
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