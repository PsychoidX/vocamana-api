package usecase

import (
	"api/model"
	"api/repository"
	"database/sql"
	"fmt"
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

func (wu *WordUsecase) GetAllWords(userId uint64) ([]model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	words, err := wu.wr.GetAllWords(userId)
	if err != nil {
		return []model.Word{}, err
	}

	return words, nil
}

func (wu *WordUsecase) GetWordById(userId uint64, wordId uint64) (model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	word, err := wu.wr.GetWordById(userId, wordId)
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
	userId := wordCreation.UserId

	// TODO: userIdがログイン中のものと一致することを確認

	createdWord, err := wu.wr.InsertWord(wordCreation)
	if err != nil {
		return model.Word{}, err
	}

	// 既存のSentenceに追加されたWord含まれればsentences_wordsに追加
	AssociateWordWithAllSentences(
		userId,
		createdWord.Id,
		wu.wr,
		wu.sr,
		wu.swr,
		wu.nr,
	)

	return createdWord, nil
}

func (wu *WordUsecase) DeleteWord(userId uint64, wordId uint64) (model.Word, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	deletedWord, err := wu.wr.DeleteWordById(userId, wordId)
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

	return updatedWord, nil
}

func (wu *WordUsecase) GetAssociatedSentencesByWordId(userId uint64, wordId uint64) ([]model.Sentence, error) {
	// TODO: userIdがログイン中のものと一致することを確認

	// wordIdの所有者がuserIdでない場合ゼロ値を返す
	isWordOwner, err := wu.wr.IsWordOwner(wordId, userId)
	if err != nil {
		return []model.Sentence{}, err
	}
	if !isWordOwner {
		return []model.Sentence{}, nil
	}

	sentences, err := wu.swr.GetAssociatedSentencesByWordId(userId, wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	// リポジトリの返り値のuserIdを検証
	userSentences := []model.Sentence{}
	for _, sentence := range sentences {
		// sentenceの所有者がuserIdでない場合continue
		isSentenceOwner, err := wu.sr.IsSentenceOwner(sentence.Id, userId)
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

func (wu *WordUsecase) GetAssociatedSentencesWithLinkByWordId(userId, wordId uint64) ([]model.SentenceWithLink, error) {
	// userIdがログイン中のものと一致するかの判定、
	// wordIdの所有者がuserIdであるかの判定は、GetAssociatedSentencesByWordId内で行われる
	userAssociatedSentences, err := wu.GetAssociatedSentencesByWordId(userId, wordId)
	if err != nil {
		return []model.SentenceWithLink{}, err
	}

	sentenceWithLinks := []model.SentenceWithLink{}
	for _, sentence := range userAssociatedSentences {
		notations, err := wu.nr.GetAllNotations(wordId)
		if err != nil {
			return []model.SentenceWithLink{}, err
		}

		// sentence.Sentence中に含まれるnotation.Notationをaタグに置換
		sentenceText := sentence.Sentence
		for _, notation := range notations {
			sentenceText = strings.Replace(
				sentenceText,
				notation.Notation,
				createWordLink(wordId, notation.Notation),
				-1,
			)
		}

		// 置換後のsentence.SentenceからSentenceWithLinkを作成
		sentenceWithLink := model.SentenceWithLink{
			Id:               sentence.Id,
			SentenceWithLink: sentenceText,
			UserId:           sentence.UserId,
			CreatedAt:        sentence.CreatedAt,
			UpdatedAt:        sentence.UpdatedAt,
		}
		sentenceWithLinks = append(sentenceWithLinks, sentenceWithLink)
	}

	return sentenceWithLinks, nil
}

func createWordLink(wordId uint64, notation string) string {
	// wordに遷移する<a>を作成
	return fmt.Sprintf(
		"<a href=\"words/%d\">%s</a>",
		wordId,
		notation,
	)
}
