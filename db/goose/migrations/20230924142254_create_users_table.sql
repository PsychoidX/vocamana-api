-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE user_id_seq;

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  email VARCHAR(400) NOT NULL UNIQUE,
  password VARCHAR(200) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER refresh_users_updated_at
  BEFORE UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE refresh_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER refresh_users_updated_at ON users;
DROP TABLE users;
DROP SEQUENCE user_id_seq;
-- +goose StatementEnd