-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS owners (
  id            UUID PRIMARY KEY,
  full_name     TEXT NOT NULL,
  phone_number  VARCHAR(500) NOT NULL,
  username      VARCHAR(500) NOT NULL,
  password_hash VARCHAR(72) NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS owners;
-- +goose StatementEnd
