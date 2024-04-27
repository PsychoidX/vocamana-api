-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE notation_id_seq;

CREATE TABLE notations (
  id INTEGER PRIMARY KEY,
  word_id INTEGER,
  notation VARCHAR(100),
  FOREIGN KEY (word_id) REFERENCES words(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER refresh_notations_updated_at
  BEFORE UPDATE ON notations FOR EACH ROW
EXECUTE PROCEDURE refresh_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER refresh_notations_updated_at ON notations;
DROP TABLE notations;
DROP SEQUENCE notation_id_seq;
-- +goose StatementEnd
