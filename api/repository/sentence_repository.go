package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type ISentenceRepository interface {
	getSequenceName() string
	InsertSentence(model.SentenceCreation) (model.Sentence, error)
}

type SentenceRepository struct {
	db *sql.DB
}

func NewSentenceRepository(db *sql.DB) ISentenceRepository {
	return &SentenceRepository{db}
}

func (sr *SentenceRepository) getSequenceName() string {
	return "sentence_id_seq"
}

func (sr *SentenceRepository) getSequenceNextvalQuery() string {
	return fmt.Sprintf("nextval('%s')", sr.getSequenceName())
}

func (sr *SentenceRepository) InsertSentence(newSentence model.SentenceCreation) (model.Sentence, error) {
	createdSentence := model.Sentence{}

	err := sr.db.QueryRow(
		"INSERT INTO sentences" +
		" (id, sentence, user_id)" +
		" VALUES(" + sr.getSequenceNextvalQuery() + ", $1, $2)" +
		" RETURNING id, sentence, user_id, created_at, updated_at;",
		newSentence.Sentence,
		newSentence.UserId,
	).Scan(
		&createdSentence.Id,
		&createdSentence.Sentence,
		&createdSentence.UserId,
		&createdSentence.CreatedAt,
		&createdSentence.UpdatedAt,
	)
	if err != nil {
		return model.Sentence{}, err
	}

	return createdSentence, nil
}