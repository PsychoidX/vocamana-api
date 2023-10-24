package repository

import (
	"api/model"
	"database/sql"
	"fmt"
)

type IWordRepository interface {
	GetAllWords(userId uint64) ([]model.Word, error)
	GetWordById(userId, wordId uint64) (model.Word, error)
	InsertWord(model.WordCreation) (model.Word, error)
	DeleteWordById(id uint64) (model.Word, error)
	UpdateWord(model.WordUpdate) (model.Word, error)
	IsWordOwner(uint64, uint64) (bool, error)
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

func (wr *WordRepository) GetAllWords(userId uint64) ([]model.Word, error) {
	var words []model.Word

	rows, err := wr.db.Query(
		"SELECT id, word, memo, user_id, created_at, updated_at FROM words" +
		" WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return []model.Word{}, err
	}
	defer rows.Close()

	for rows.Next() {
		word := model.Word{}
		err:= rows.Scan(&word.Id, &word.Word, &word.Memo, &word.UserId, &word.CreatedAt, &word.UpdatedAt);
		if err != nil {
			return []model.Word{}, err
		}
		words = append(words, word)
	}

	return words, nil
}

func (wr *WordRepository) GetWordById(userId uint64, wordId uint64) (model.Word, error) {
	word := model.Word{}

	err := wr.db.QueryRow(
		"SELECT id, word, memo, user_id, created_at, updated_at" + 
		" FROM words" +
		" WHERE id = $1" +
		" AND user_id = $2;",
		wordId,
		userId,
	).Scan(&word.Id, &word.Word, &word.Memo, &word.UserId, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		return model.Word{}, err
	}

	return word, nil
}

func (wr *WordRepository) InsertWord(newWord model.WordCreation) (model.Word, error) {
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

func (wr *WordRepository) DeleteWordById(id uint64) (model.Word, error) {
	deletedWord := model.Word{}

	err := wr.db.QueryRow(
		"DELETE FROM words" +
		" WHERE id = $1" +
		" RETURNING id, word, memo, user_id, created_at, updated_at;",
		id,
	).Scan(
		&deletedWord.Id,
		&deletedWord.Word,
		&deletedWord.Memo,
		&deletedWord.UserId,
		&deletedWord.CreatedAt,
		&deletedWord.UpdatedAt,
	)
	if err != nil {
		return model.Word{}, err
	}

	return deletedWord, nil
}

func (wr *WordRepository) UpdateWord(wordUpdate model.WordUpdate) (model.Word, error) {
	updatedWord := model.Word{}

	err := wr.db.QueryRow(
		"UPDATE words" +
		" SET word = $1," +
		" memo = $2" +
		" WHERE id = $3" +
		" RETURNING id, word, memo, user_id, created_at, updated_at;",
		wordUpdate.Word,
		wordUpdate.Memo,
		wordUpdate.Id,
	).Scan(
		&updatedWord.Id,
		&updatedWord.Word,
		&updatedWord.Memo,
		&updatedWord.UserId,
		&updatedWord.CreatedAt,
		&updatedWord.UpdatedAt,
	)
	if err != nil {
		return model.Word{}, err
	}

	return updatedWord, nil
}

func (wr *WordRepository) IsWordOwner(wordId uint64, userId uint64) (bool, error) {
	// wordIdの所持者がuserIdであるかを判定

	var count int

	err := wr.db.QueryRow(
		"SELECT COUNT(*) FROM words" +
		" WHERE id = $1" + 
		" AND user_id = $2;",
		wordId,
		userId,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count == 1, nil
}