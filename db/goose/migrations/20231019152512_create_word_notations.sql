-- +goose Up
-- +goose StatementBegin
CREATE TABLE word_notations (
  word_id INTEGER PRIMARY KEY,
  notation VARCHAR(100),
  FOREIGN KEY (word_id) REFERENCES words(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE word_notations;
-- +goose StatementEnd
