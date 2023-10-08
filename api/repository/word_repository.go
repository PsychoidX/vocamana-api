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
	var id uint

	err := wr.db.QueryRow(
		"INSERT INTO words" +
		" (id, word, memo)" +
		" VALUES(" + wr.getSequenceNextvalQuery() + ", $1, $2)" +
		" RETURNING id;",
		newWord.Word,
		newWord.Memo,
	).Scan(&id)
	if err != nil {
		return model.Word{}, err
	}

	return wr.GetWordById(id)
}