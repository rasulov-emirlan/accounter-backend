-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS categories (
  id                 uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  store_id           uuid NOT NULL,
  parent_category_id uuid,
  name               VARCHAR(255) NOT NULL,
  article            VARCHAR(100),
  icon_url           TEXT,
  created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_categories_store_id FOREIGN KEY (store_id)
    REFERENCES stores(id),
  CONSTRAINT fk_categories_parent_category_id FOREIGN KEY (parent_category_id)
    REFERENCES categories(id),
  CONSTRAINT check_categories_parent_loop CHECK (
    parent_category_id IS NULL OR
    id <> parent_category_id
  )
);

CREATE INDEX IF NOT EXISTS ix_categories_name_trgm ON categories USING GIN(name gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
DROP INDEX IF EXISTS ix_categories_name_trgm;
DROP EXTENSION IF EXISTS "pg_trgm";
-- +goose StatementEnd
