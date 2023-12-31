package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type ISentenceRepository interface {
	GetAllSentences(userId uint64) ([]model.Sentence, error)
	GetAllSentencesWithLimit(userId uint64, limit uint64, offset uint64) ([]model.Sentence, error)
	GetSentenceById(userId uint64, sentenceId uint64) (model.Sentence, error)
	InsertSentence(model.SentenceCreation) (model.Sentence, error)
	UpdateSentence(model.SentenceUpdate) (model.Sentence, error)
	DeleteSentenceById(userId uint64, sentenceId uint64) (model.Sentence, error)
	IsSentenceOwner(sentenceId uint64, userId uint64) (bool, error)
	GetSentencesCount(userId uint64) (uint64, error)
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

	rows, err := sr.db.Query(`
		SELECT id, sentence, user_id, created_at, updated_at FROM sentences
		WHERE user_id = $1;
		`,
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

func (sr *SentenceRepository) GetAllSentencesWithLimit(userId, limit, offset uint64) ([]model.Sentence, error) {
	var sentences []model.Sentence

	rows, err := sr.db.Query(`
		SELECT id, sentence, user_id, created_at, updated_at FROM sentences
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2
		OFFSET $3;
		`,
		userId,
		limit,
		offset,
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

func (sr *SentenceRepository) GetSentenceById(userId, sentenceId uint64) (model.Sentence, error) {
	sentence := model.Sentence{}

	err := sr.db.QueryRow(`
		SELECT id, sentence, user_id, created_at, updated_at
		FROM sentences
		WHERE id = $1
			AND user_id = $2
		`,
		sentenceId,
		userId,
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
		newSentence.LoginUserId,
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

func (sr *SentenceRepository) UpdateSentence(sentenceUpdate model.SentenceUpdate) (model.Sentence, error) {
	updatedSentence := model.Sentence{}

	err := sr.db.QueryRow(`
		UPDATE sentences
		SET sentence = $1
		WHERE user_id = $2
			AND id = $3
		RETURNING id, sentence, user_id, created_at, updated_at;
		`,
		sentenceUpdate.Sentence,
		sentenceUpdate.LoginUserId,
		sentenceUpdate.Id,
	).Scan(
		&updatedSentence.Id,
		&updatedSentence.Sentence,
		&updatedSentence.UserId,
		&updatedSentence.CreatedAt,
		&updatedSentence.UpdatedAt,
	)
	if err != nil {
		return model.Sentence{}, err
	}

	return updatedSentence, nil
}

func (sr *SentenceRepository) DeleteSentenceById(userId, sentenceId uint64) (model.Sentence, error) {
	deletedSentence := model.Sentence{}

	err := sr.db.QueryRow(`
		DELETE FROM sentences
		WHERE user_id = $1
			AND id = $2
		RETURNING id, sentence, user_id, created_at, updated_at;
		`,
		userId,
		sentenceId,
	).Scan(
		&deletedSentence.Id,
		&deletedSentence.Sentence,
		&deletedSentence.UserId,
		&deletedSentence.CreatedAt,
		&deletedSentence.UpdatedAt,
	)
	if err != nil {
		return model.Sentence{}, err
	}

	return deletedSentence, nil
}

func (sr *SentenceRepository) IsSentenceOwner(sentenceId uint64, userId uint64) (bool, error) {
	// sentenceIdの所持者がuserIdであるかを判定
	var count int
	err := sr.db.QueryRow(
		"SELECT COUNT(*) FROM sentences" +
		" WHERE id = $1" + 
		" AND user_id = $2;",
		sentenceId,
		userId,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count == 1, nil
}

func (sr *SentenceRepository) GetSentencesCount(userId uint64) (uint64, error) {
	var count uint64

	err := sr.db.QueryRow(`
		SELECT COUNT(*)
		FROM sentences
		WHERE user_id = $1
		`,
		userId,
	).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}