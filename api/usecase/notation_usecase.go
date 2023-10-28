package usecase

import (
	"api/model"
	"api/repository"
)

type NotationUsecase struct {
	nr repository.INotationRepository
	wr repository.IWordRepository
}

func NewNotationUsecase(
		nr repository.INotationRepository,
		wr repository.IWordRepository,
	) *NotationUsecase {
	return &NotationUsecase{nr, wr}
}

func (nu *NotationUsecase) GetAllNotations(userId, wordId uint64) ([]model.Notation, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	// wordIdの所有者がuserIdの場合ゼロ値を返す
	isWordOwner, err := nu.wr.IsWordOwner(wordId, userId)
	if err != nil {
		return []model.Notation{}, err
	}
	if !isWordOwner {
		return []model.Notation{}, nil
	}

	notations, err := nu.nr.GetAllNotations(wordId)
	if err != nil {
		return []model.Notation{}, err
	}

	return notations, nil
}

func (nu *NotationUsecase) CreateNotation(userId uint64, notationCreation model.NotationCreation) (model.Notation, error) {
	// TODO: userIdがログイン中のものと一致することを確認
	
	// 追加先のWordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := nu.wr.IsWordOwner(notationCreation.WordId, userId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}
	
	createdNotation, err := nu.nr.InsertNotation(notationCreation)
	if err != nil {
		return model.Notation{}, err
	}

	return createdNotation, nil
}