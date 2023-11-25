package repository

import (
	"api/model"
	"database/sql"
)

type ISentencesWordsRepository interface {
	AssociateSentenceWithWord(sentenceId uint64, wordId uint64) error
	GetUserAssociatedSentencesByWordId(wordId uint64) ([]model.Sentence, error)
	GetUserAssociatedWordsBySentenceId(sentenceId uint64) ([]model.Word, error)
	DeleteAllAssociationBySentenceId(sentenceId uint64) error
	DeleteAllAssociationByWordId(sentenceId uint64) error
}

type SentencesWordsRepository struct {
	db *sql.DB
}

func NewSentencesWordsRepository(db *sql.DB) ISentencesWordsRepository {
	return &SentencesWordsRepository{db}
}

func (swr *SentencesWordsRepository) AssociateSentenceWithWord(sentenceId uint64, wordId uint64) error {
	// テーブルにレコードが存在しない場合にだけ追加
	
	// sentenceIdのSentenceとwordIdのWordのuserIdは一致する（SentenceとWordの所有者は同じである）ことを前提とする
	// AssociateSentenceWithWord() 呼び出し時に常軌を前提とするため、
	// GetAssociatedSentencesByWordId(), GetAssociatedWordsBySentenceId() 等の取得メソッドではuserIdの検証を行わない

	_, err := swr.db.Exec(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		SELECT $1, $2
		WHERE NOT EXISTS(
			SELECT 1
			FROM sentences_words
			WHERE sentence_id = $1 AND word_id = $2
		);`,
		sentenceId,
		wordId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (swr *SentencesWordsRepository) GetUserAssociatedSentencesByWordId(wordId uint64) ([]model.Sentence, error) {
	// wordIdに紐づくSentenceを全件取得
	// sentenceIdのSentenceとwordIdのWordのuserIdは一致する（SentenceとWordの所有者は同じである）ことを前提とするため、
	// 参照時にはuserIdの検証を行わない

	var sentences []model.Sentence

	rows, err := swr.db.Query(`
		SELECT
			sentences.id,
			sentences.sentence,
			sentences.user_id,
			sentences.created_at,
			sentences.updated_at
		FROM sentences_words
		LEFT JOIN words
			ON sentences_words.word_id = words.id
		LEFT JOIN sentences
			ON sentences_words.sentence_id = sentences.id
		WHERE sentences_words.word_id = $1;
		`,
		wordId,
	)
	if err != nil {
		return []model.Sentence{}, err
	}
	defer rows.Close()

	for rows.Next() {
		sentence := model.Sentence{}
		err := rows.Scan(&sentence.Id, &sentence.Sentence, &sentence.UserId, &sentence.CreatedAt, &sentence.UpdatedAt)
		if err != nil {
			return []model.Sentence{}, err
		}
		sentences = append(sentences, sentence)
	}

	return sentences, nil
}

func (swr *SentencesWordsRepository) GetUserAssociatedWordsBySentenceId(sentenceId uint64) ([]model.Word, error) {
	// sentenceIdに紐づくWordのを全件取得
	// sentenceIdのSentenceとwordIdのWordのuserIdは一致する（SentenceとWordの所有者は同じである）ことを前提とするため、
	// 参照時にはuserIdの検証を行わない
	
	var words []model.Word

	rows, err := swr.db.Query(`
		SELECT
			words.id,
			words.word,
			words.memo,
			words.user_id,
			words.created_at,
			words.updated_at
		FROM sentences_words
		LEFT JOIN words
			ON sentences_words.word_id = words.id
		LEFT JOIN sentences
			ON sentences_words.sentence_id = sentences.id
		WHERE sentences_words.sentence_id = $1
		`,
		sentenceId,
	)
	if err != nil {
		return []model.Word{}, err
	}
	defer rows.Close()

	for rows.Next() {
		word := model.Word{}
		err := rows.Scan(&word.Id, &word.Word, &word.Memo, &word.UserId, &word.CreatedAt, &word.UpdatedAt)
		if err != nil {
			return []model.Word{}, err
		}
		words = append(words, word)
	}

	return words, nil
}

func (swr *SentencesWordsRepository) DeleteAllAssociationBySentenceId(sentenceId uint64) error {
	_, err := swr.db.Exec(`
		DELETE FROM sentences_words
		WHERE sentence_id = $1;
		`,
		sentenceId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (swr *SentencesWordsRepository) DeleteAllAssociationByWordId(wordId uint64) error {
	_, err := swr.db.Exec(`
		DELETE FROM sentences_words
		WHERE word_id = $1;
		`,
		wordId,
	)
	if err != nil {
		return err
	}

	return nil
}