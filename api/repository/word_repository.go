package repository

import (
	"api/model"
	"database/sql"
)

type IWordRepository interface {
	GetAllWords() ([]model.Word, error)
	GetWordById(id uint) (model.Word, error)
}

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) IWordRepository {
	return &WordRepository{db}
}

func (wr *WordRepository) GetAllWords() ([]model.Word, error) {
	// TODO
	return []model.Word{}, nil
}

func (wr *WordRepository) GetWordById(id uint) (model.Word, error) {
	// row := wr.db.QueryRow("SELECT id, word, memo, created_at, updated_at FROM WORDS")
	// err := row.Scan(&word.Id, &word.Memo, &word.CreatedAt, &word.UpdatedAt)
	// if err != nil {
	// 	return nil, err
	// }
	return model.Word{}, nil
}