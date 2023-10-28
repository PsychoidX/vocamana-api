package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type INotationRepository interface {
	InsertNotation(model.NotationCreation) (model.Notation, error)
}

type NotationRepository struct {
	db *sql.DB
}

func NewNotationRepository(db *sql.DB) INotationRepository {
	return &NotationRepository{db}
}

func (nr *NotationRepository) getSequenceName() string {
	return "notation_id_seq"
}

func (nr *NotationRepository) getSequenceNextvalQuery() string {
	return fmt.Sprintf("nextval('%s')", nr.getSequenceName())
}

func (nr *NotationRepository) InsertNotation(notationCreation model.NotationCreation) (model.Notation, error) {
	createdNotation := model.Notation{}

	err := nr.db.QueryRow(fmt.Sprintf(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(%s, $1, $2)
		RETURNING id, word_id, notation, created_at, updated_at;
		`, 
		nr.getSequenceNextvalQuery(),
		),
		notationCreation.WordId,
		notationCreation.Notation,
	).Scan(
		&createdNotation.Id,
		&createdNotation.WordId,
		&createdNotation.Notation,
		&createdNotation.CreatedAt,
		&createdNotation.UpdatedAt,
	)
	if err != nil {
		return model.Notation{}, err
	}

	return createdNotation, nil
}