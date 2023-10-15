-- +goose Up
-- +goose StatementBegin
CREATE TABLE containing (
  sentence_id INTEGER,
  word_id INTEGER,
  PRIMARY KEY(sentence_id, word_id),
  FOREIGN KEY (sentence_id) REFERENCES sentences(id) ON DELETE CASCADE,
  FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE containing;
-- +goose StatementEnd
