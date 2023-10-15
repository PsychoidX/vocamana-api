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