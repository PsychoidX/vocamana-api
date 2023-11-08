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

func (wu *WordUsecase) GetAllWords(loginUserId uint64) ([]model.Word, error) {
	words, err := wu.wr.GetAllWords(loginUserId)
	if err != nil {
		return []model.Word{}, err
	}

	return words, nil
}

func (wu *WordUsecase) GetWordById(loginUserId, wordId uint64) (model.Word, error) {
	word, err := wu.wr.GetWordById(loginUserId, wordId)
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
	loginUserId := wordCreation.UserId

	createdWord, err := wu.wr.InsertWord(wordCreation)
	if err != nil {
		return model.Word{}, err
	}

	// 既存のSentenceに追加されたWord含まれればsentences_wordsに追加
	AssociateWordWithAllSentences(
		loginUserId,
		createdWord.Id,
		wu.wr,
		wu.sr,
		wu.swr,
		wu.nr,
	)

	return createdWord, nil
}

func (wu *WordUsecase) CreateMultipleWord(wordCreations []model.WordCreation) ([]model.Word, error) {
	// TODO 1件でも失敗したらロールバックする実装に変更
	var createdWords []model.Word
	for _, wordCreation := range wordCreations {
		createdWord, err := wu.CreateWord(wordCreation)
		if err != nil {
			return []model.Word{}, err
		}

		createdWords = append(createdWords, createdWord)
	}

	return createdWords, nil
}

func (wu *WordUsecase) DeleteWord(loginUserId, wordId uint64) (model.Word, error) {
	deletedWord, err := wu.wr.DeleteWordById(loginUserId, wordId)
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

func (wu *WordUsecase) GetAssociatedSentencesByWordId(loginUserId, wordId uint64) ([]model.Sentence, error) {
	// wordIdの所有者がloginUserIdでない場合ゼロ値を返す
	isWordOwner, err := wu.wr.IsWordOwner(wordId, loginUserId)
	if err != nil {
		return []model.Sentence{}, err
	}
	if !isWordOwner {
		return []model.Sentence{}, nil
	}

	sentences, err := wu.swr.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	// リポジトリの返り値のuserIdを検証
	userSentences := []model.Sentence{}
	for _, sentence := range sentences {
		// sentenceの所有者がloginUserIdでない場合continue
		isSentenceOwner, err := wu.sr.IsSentenceOwner(sentence.Id, loginUserId)
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

func (wu *WordUsecase) GetAssociatedSentencesWithLinkByWordId(loginUserId, wordId uint64) ([]model.SentenceWithLink, error) {
	userAssociatedSentences, err := wu.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return []model.SentenceWithLink{}, err
	}

	sentenceWithLinks := []model.SentenceWithLink{}
	for _, sentence := range userAssociatedSentences {
		sentenceText := sentence.Sentence

		// sentenceに紐づくWordを全件取得し、sentence中におけるそのWordの出現箇所をリンクに変換
		words, err := wu.swr.GetAssociatedWordsBySentenceId(loginUserId, sentence.Id)
		if err != nil {
			return []model.SentenceWithLink{}, err
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
			notations, err := wu.nr.GetAllNotations(word.Id)
			if err != nil {
				return []model.SentenceWithLink{}, err
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
		"<a href=\"/words/%d\">%s</a>",
		wordId,
		notation,
	)
}
