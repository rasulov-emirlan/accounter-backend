-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stores (
  id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner_id    uuid NOT NULL,
  name        TEXT NOT NULL,
  description TEXT NOT NULL,
  tsv         TSVECTOR,
  created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT  fk_stores_owner_id FOREIGN KEY (owner_id)
    REFERENCES owners(id)
);

CREATE INDEX IF NOT EXISTS ix_stores_tsv ON stores USING GIN(tsv);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
DROP INDEX IF EXISTS ix_stores_tsv;
-- +goose StatementEnd
