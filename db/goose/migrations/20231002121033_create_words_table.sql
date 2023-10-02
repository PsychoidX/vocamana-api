-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE word_id_seq;

CREATE TABLE words (
  id INTEGER PRIMARY KEY,
  word VARCHAR(100) NOT NULL,
  memo VARCHAR(500),
  user_id INTEGER,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TRIGGER refresh_words_updated_at
  BEFORE UPDATE ON words FOR EACH ROW
EXECUTE PROCEDURE refresh_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER refresh_words_updated_at ON words;
DROP TABLE words;
DROP SEQUENCE word_id_seq;
-- +goose StatementEnd
