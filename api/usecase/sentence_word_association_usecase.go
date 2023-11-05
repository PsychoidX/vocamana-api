package usecase

import (
	"api/model"
	"api/repository"
	"strings"
)

func AssociateWordWithAllSentences(
	userId uint64,
	wordId uint64,
	wr repository.IWordRepository,
	sr repository.ISentenceRepository,
	swr repository.ISentencesWordsRepository,
	nr repository.INotationRepository,
) ([]model.Sentence, error) {
	// userIdに紐づく全Sentenceに対し、
	// Sentenceの中にwordIdのWordまたはNotationが含まれればsentences_wordsにレコード追加

	// TODO: userIdがログイン中のものと一致することを確認

	// wordIdの所有者がuserIdでない場合何もしない
	isWordOwner, err := wr.IsWordOwner(wordId, userId)
	if err != nil {
		return []model.Sentence{}, err
	}
	if !isWordOwner {
		return []model.Sentence{}, nil
	}

	word, err := wr.GetWordById(userId, wordId)
	if err != nil {
		return []model.Sentence{}, err
	}

	userSentences, err := sr.GetAllSentences(userId)
	if err != nil {
		return []model.Sentence{}, err
	}

	var associatedSentences []model.Sentence
	for _, sentence := range userSentences {
		// Sentence中にWordが含まれるか判定
		if strings.Contains(sentence.Sentence, word.Word) {
			err = swr.AssociateSentenceWithWord(sentence.Id, word.Id)
			if err != nil {
				return []model.Sentence{}, err
			}
			associatedSentences = append(associatedSentences, sentence)
			// Sentenceの中にWordが含まれる場合、
			// continueし、Notationが含まれるかの判定はしない
			continue
		}

		// Sentence中にNotationが含まれるか判定
		notations, err := nr.GetAllNotations(word.Id)
		if err != nil {
			return []model.Sentence{}, err
		}

		for _, notation := range notations {
			if strings.Contains(sentence.Sentence, notation.Notation) {
				err = swr.AssociateSentenceWithWord(sentence.Id, word.Id)
				if err != nil {
					return []model.Sentence{}, err
				}
				associatedSentences = append(associatedSentences, sentence)
				// Sentenceの中にNotationが含まれる場合、
				// sentences_wordsに2つ目のレコードが追加されないようbreak
				break
			}
		}
	}

	return associatedSentences, nil
}
