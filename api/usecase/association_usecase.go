package usecase

import (
	"api/model"
	"api/repository"
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

	associatedUserSentences, err := au.swr.GetUserAssociatedSentencesByWordId(wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	return associatedUserSentences, nil
}

func (au *AssociationUsecase) GetAssociatedSentencesWithLinkByWordId(loginUserId, wordId uint64) ([]model.SentenceWithLink, error) {
	userAssociatedSentences, err := au.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return []model.SentenceWithLink{}, err
	}

	sentenceWithLinks := []model.SentenceWithLink{}
	for _, sentence := range userAssociatedSentences {
		sentenceWithLink, err := au.toSentenceWithLink(loginUserId, sentence)
		if err != nil {
			return []model.SentenceWithLink{}, err
		}

		sentenceWithLinks = append(sentenceWithLinks, sentenceWithLink)
	}

	return sentenceWithLinks, nil
}

func (au *AssociationUsecase) GetSentenceWithLinkById(loginUserId, sentenceId uint64) (model.SentenceWithLink, error) {
	sentence, err := au.su.GetSentenceById(loginUserId, sentenceId)
	if err != nil {
		return model.SentenceWithLink{}, err
	}

	sentenceWithLink, err := au.toSentenceWithLink(loginUserId, sentence)
	if err != nil {
		return model.SentenceWithLink{}, err
	}

	return sentenceWithLink, nil
}

func (au *AssociationUsecase) GetAllSentencesWithLink(loginUserId uint64) ([]model.SentenceWithLink, error) {
	sentences, err := au.su.GetAllSentences(loginUserId)
	if err != nil {
		return []model.SentenceWithLink{}, err
	}

	sentenceWithLinks := []model.SentenceWithLink{}
	for _, sentence := range sentences {
		sentenceWithLink, err := au.toSentenceWithLink(loginUserId, sentence)
		if err != nil {
			return []model.SentenceWithLink{}, err
		}

		sentenceWithLinks = append(sentenceWithLinks, sentenceWithLink)
	}

	return sentenceWithLinks, nil
}

func (au *AssociationUsecase) toSentenceWithLink(loginUserId uint64, sentence model.Sentence) (model.SentenceWithLink, error) {
	// sentenceに紐づくWordを全件取得し、sentence中におけるそのWordの出現箇所をリンクに変換
	sentenceText := sentence.Sentence

	words, err := au.swr.GetUserAssociatedWordsBySentenceId(sentence.Id)
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