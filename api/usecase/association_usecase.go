package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
	"fmt"
	"strings"
)

type AssociationUsecase struct {
	wr  repository.IWordRepository
	sr  repository.ISentenceRepository
	swr repository.ISentencesWordsRepository
	nr  repository.INotationRepository
	wu  *WordUsecase
	su  *SentenceUsecase
}

func NewAssociationUsecase(
	wr repository.IWordRepository,
	sr repository.ISentenceRepository,
	swr repository.ISentencesWordsRepository,
	nr repository.INotationRepository,
) *AssociationUsecase {
	wu := NewWordUsecase(wr, sr, swr, nr)
	su := NewSentenceUsecase(sr, wr, swr, nr)
	return &AssociationUsecase{wr, sr, swr, nr, wu, su}
}

func (au *AssociationUsecase) GetAssociatedSentencesByWordId(loginUserId, wordId uint64) ([]model.Sentence, error) {
	// wordIdの所有者がloginUserIdでない場合ゼロ値を返す
	isWordOwner, err := au.wr.IsWordOwner(wordId, loginUserId)
	if err != nil {
		return []model.Sentence{}, err
	}
	if !isWordOwner {
		return []model.Sentence{}, nil
	}

	sentences, err := au.swr.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	// リポジトリの返り値のuserIdを検証
	userSentences := []model.Sentence{}
	for _, sentence := range sentences {
		// sentenceの所有者がloginUserIdでない場合continue
		isSentenceOwner, err := au.sr.IsSentenceOwner(sentence.Id, loginUserId)
		if err != nil {
			return []model.Sentence{}, err
		}
		if !isSentenceOwner {
			continue
		}

		userSentences = append(userSentences, sentence)
	}

	return userSentences, nil
}

func (au *AssociationUsecase) GetAssociatedSentencesWithLinkByWordId(loginUserId, wordId uint64) ([]model.SentenceWithLink, error) {
	userAssociatedSentences, err := au.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return []model.SentenceWithLink{}, err
	}

	sentenceWithLinks := []model.SentenceWithLink{}
	for _, sentence := range userAssociatedSentences {
		sentenceWithLink, err := au.ToSentenceWithLink(loginUserId, sentence)
		if err != nil {
			return []model.SentenceWithLink{}, err
		}

		sentenceWithLinks = append(sentenceWithLinks, sentenceWithLink)
	}

	return sentenceWithLinks, nil
}

func (au *AssociationUsecase) GetSentencesWithLinkById(loginUserId, sentenceId uint64) (model.SentenceWithLink, error) {
	sentence, err := au.sr.GetSentenceById(loginUserId, sentenceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// マッチするレコードが無い場合
			// Sentenceのゼロ値を返す
			return model.SentenceWithLink{}, nil
		}

		return model.SentenceWithLink{}, err
	}

	sentenceWithLink, err := au.ToSentenceWithLink(loginUserId, sentence)
	if err != nil {
		return model.SentenceWithLink{}, err
	}

	return sentenceWithLink, nil
}

func (au *AssociationUsecase) ToSentenceWithLink(loginUserId uint64, sentence model.Sentence) (model.SentenceWithLink, error) {
	// sentenceに紐づくWordを全件取得し、sentence中におけるそのWordの出現箇所をリンクに変換
	sentenceText := sentence.Sentence

	words, err := au.swr.GetAssociatedWordsBySentenceId(loginUserId, sentence.Id)
	if err != nil {
		return model.SentenceWithLink{}, err
	}

	for _, word := range words {
		// sentenceText中に含まれるword.Wordをaタグに置換
		sentenceText = strings.Replace(
			sentenceText,
			word.Word,
			createWordLink(word.Id, word.Word),
			-1,
		)

		// sentenceText中に含まれるnotation.Notationをaタグに置換
		notations, err := au.nr.GetAllNotations(word.Id)
		if err != nil {
			return model.SentenceWithLink{}, err
		}

		for _, notation := range notations {
			sentenceText = strings.Replace(
				sentenceText,
				notation.Notation,
				createWordLink(word.Id, notation.Notation),
				-1,
			)
		}
	}

	// 置換後のsentenceTextからSentenceWithLinkを作成
	sentenceWithLink := model.SentenceWithLink{
		Id:               sentence.Id,
		Sentence:         sentence.Sentence,
		SentenceWithLink: sentenceText,
		UserId:           sentence.UserId,
		CreatedAt:        sentence.CreatedAt,
		UpdatedAt:        sentence.UpdatedAt,
	}

	return sentenceWithLink, nil
}

func createWordLink(wordId uint64, notation string) string {
	// wordに遷移する<a>を作成
	return fmt.Sprintf(
		"<a href=\"/words/%d\">%s</a>",
		wordId,
		notation,
	)
}