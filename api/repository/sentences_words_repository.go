package repository

import (
	"database/sql"
)

type ISentencesWordsRepository interface {
	AssociateSentenceWithWord(uint64, uint64) error
}

type SentencesWordsRepository struct {
	db *sql.DB
}

func NewSentencesWordsRepository(db *sql.DB) ISentencesWordsRepository {
	return &SentencesWordsRepository{db}
}

func (swr *SentencesWordsRepository) AssociateSentenceWithWord(sentenceId uint64, wordId uint64) (error) {
	// テーブルにレコードが存在しない場合にだけ追加
	_, err := swr.db.Exec(
		"INSERT INTO sentences_words" +
		" (sentence_id, word_id)" +
		" SELECT $1, $2" +
		" WHERE NOT EXISTS(" +
		"   SELECT 1" +
		"   FROM sentences_words" +
		"   WHERE sentence_id = $1 AND word_id = $2" +
		" );",
		sentenceId,
		wordId,
	)
	if err != nil {
		return err
	}

	return nil
}