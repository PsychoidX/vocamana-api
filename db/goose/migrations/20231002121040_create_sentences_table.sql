-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE sentence_id_seq;

CREATE TABLE sentences (
  id INTEGER PRIMARY KEY,
  sentence VARCHAR(500) NOT NULL,
  user_id INTEGER,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TRIGGER refresh_sentences_updated_at
  BEFORE UPDATE ON sentences FOR EACH ROW
EXECUTE PROCEDURE refresh_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER refresh_sentences_updated_at ON sentences;
DROP TABLE sentences;
DROP SEQUENCE sentence_id_seq;
-- +goose StatementEnd
