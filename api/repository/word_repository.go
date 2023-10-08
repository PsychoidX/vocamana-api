package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type IWordRepository interface {
	getSequenceName() string
	GetAllWords() ([]model.Word, error)
	GetWordById(id uint) (model.Word, error)
	InsertWord(model.WordRegistration) (model.Word, error)
}

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) IWordRepository {
	return &WordRepository{db}
}

func (wr *WordRepository) getSequenceName() string {
	return "word_id_seq"
}

func (wr *WordRepository) getSequenceNextvalQuery() string {
	return fmt.Sprintf("nextval('%s')", wr.getSequenceName())
}

func (wr *WordRepository) GetAllWords() ([]model.Word, error) {
	// TODO
	return []model.Word{}, nil
}

func (wr *WordRepository) GetWordById(id uint) (model.Word, error) {
	word := model.Word{}

	err := wr.db.QueryRow(
		"SELECT id, word, memo, created_at, updated_at" + 
		" FROM words" +
		" WHERE id = $1",
		id,
	).Scan(&word.Id, &word.Word, &word.Memo, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		return model.Word{}, err
	}

	return word, nil
}

func (wr *WordRepository) InsertWord(newWord model.WordRegistration) (model.Word, error) {
	createdWord := model.Word{}

	err := wr.db.QueryRow(
		"INSERT INTO words" +
		" (id, word, memo, user_id)" +
		" VALUES(" + wr.getSequenceNextvalQuery() + ", $1, $2, $3)" +
		" RETURNING id, word, memo, user_id, created_at, updated_at;",
		newWord.Word,
		newWord.Memo,
		newWord.UserId,
	).Scan(
		&createdWord.Id,
		&createdWord.Word,
		&createdWord.Memo,
		&createdWord.UserId,
		&createdWord.CreatedAt,
		&createdWord.UpdatedAt,
	)
	if err != nil {
		return model.Word{}, err
	}

	return createdWord, nil
}