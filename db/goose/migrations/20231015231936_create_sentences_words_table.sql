-- +goose Up
-- +goose StatementBegin
CREATE TABLE sentences_words (
  sentence_id INTEGER,
  word_id INTEGER,
  PRIMARY KEY(sentence_id, word_id),
  FOREIGN KEY (sentence_id) REFERENCES sentences(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  FOREIGN KEY (word_id) REFERENCES words(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sentences_words;
-- +goose StatementEnd
