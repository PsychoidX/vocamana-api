package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
)

type NotationUsecase struct {
	wr  repository.IWordRepository
	sr  repository.ISentenceRepository
	swr repository.ISentencesWordsRepository
	nr  repository.INotationRepository
}

func NewNotationUsecase(
	wr repository.IWordRepository,
	sr repository.ISentenceRepository,
	swr repository.ISentencesWordsRepository,
	nr repository.INotationRepository,
) *NotationUsecase {
	return &NotationUsecase{wr, sr, swr, nr}
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

	// 既存のSentenceに追加されたWord含まれればsentences_wordsに追加
	AssociateWordWithAllSentences(
		userId,
		createdNotation.WordId,
		nu.wr,
		nu.sr,
		nu.swr,
		nu.nr,
	)

	return createdNotation, nil
}

func (nu *NotationUsecase) UpdateNotation(userId uint64, notationUpdate model.NotationUpdate) (model.Notation, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	// WordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := nu.wr.IsWordOwner(notationUpdate.WordId, userId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}

	updatedNotation, err := nu.nr.UpdateNotation(notationUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	return updatedNotation, nil
}

func (nu *NotationUsecase) DeleteNotation(userId, wordId, notationId uint64) (model.Notation, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	// WordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := nu.wr.IsWordOwner(wordId, userId)
	if err != nil {
		return model.Notation{}, err
	}
	if !isWordOwner {
		return model.Notation{}, nil
	}

	deletedNotation, err := nu.nr.DeleteNotationById(notationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// レコードが更新されなかった場合
			// Notationのゼロ値を返す
			return model.Notation{}, nil
		}

		return model.Notation{}, err
	}

	return deletedNotation, nil
}
