-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION refresh_updated_at() RETURNS trigger AS
$$
BEGIN
  NEW.updated_at := CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION refresh_updated_at();
-- +goose StatementEnd
