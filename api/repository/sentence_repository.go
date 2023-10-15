package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type ISentenceRepository interface {
	getSequenceName() string
	GetAllSentences(uint64) ([]model.Sentence, error)
	GetSentenceById(uint64) (model.Sentence, error)
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

func (sr *SentenceRepository) GetAllSentences(userId uint64) ([]model.Sentence, error) {
	var sentences []model.Sentence

	rows, err := sr.db.Query(
		"SELECT id, sentence, user_id, created_at, updated_at FROM sentences" +
		" WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return []model.Sentence{}, err
	}
	defer rows.Close()

	for rows.Next() {
		sentence := model.Sentence{}
		err:= rows.Scan(&sentence.Id, &sentence.Sentence, &sentence.UserId, &sentence.CreatedAt, &sentence.UpdatedAt);
		if err != nil {
			return []model.Sentence{}, err
		}
		sentences = append(sentences, sentence)
	}

	return sentences, nil
}

func (sr *SentenceRepository) GetSentenceById(id uint64) (model.Sentence, error) {
	sentence := model.Sentence{}

	err := sr.db.QueryRow(
		"SELECT id, sentence, user_id, created_at, updated_at" + 
		" FROM sentences" +
		" WHERE id = $1",
		id,
	).Scan(&sentence.Id, &sentence.Sentence, &sentence.UserId, &sentence.CreatedAt, &sentence.UpdatedAt)
	if err != nil {
		return model.Sentence{}, err
	}

	return sentence, nil
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