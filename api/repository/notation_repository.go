package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type INotationRepository interface {
	GetAllNotations(uint64) ([]model.Notation, error)
	InsertNotation(model.NotationCreation) (model.Notation, error)
	UpdateNotation(model.NotationUpdate) (model.Notation, error)
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

func (nr *NotationRepository) GetAllNotations(wordId uint64) ([]model.Notation, error) {
	var notations []model.Notation

	rows, err := nr.db.Query(`
		SELECT id, word_id, notation, created_at, updated_at FROM notations
		WHERE word_id = $1
		`,
		wordId,
	)
	if err != nil {
		return []model.Notation{}, err
	}
	defer rows.Close()

	for rows.Next() {
		notation := model.Notation{}
		err := rows.Scan(
			&notation.Id,
			&notation.WordId,
			&notation.Notation,
			&notation.CreatedAt,
			&notation.UpdatedAt,
		);
		if err != nil {
			return []model.Notation{}, err
		}
		notations = append(notations, notation)
	}

	return notations, nil
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

func (nr *NotationRepository) UpdateNotation(notationUpdate model.NotationUpdate) (model.Notation, error) {
	updatedNotation := model.Notation{}

	err := nr.db.QueryRow(`
		UPDATE notations
		SET notation = $1
		WHERE word_id = $2
			AND id = $3
		RETURNING id, word_id, notation, created_at, updated_at;
		`, 
		notationUpdate.Notation,
		notationUpdate.WordId,
		notationUpdate.Id,
	).Scan(
		&updatedNotation.Id,
		&updatedNotation.WordId,
		&updatedNotation.Notation,
		&updatedNotation.CreatedAt,
		&updatedNotation.UpdatedAt,
	)
	if err != nil {
		return model.Notation{}, err
	}

	return updatedNotation, nil
}