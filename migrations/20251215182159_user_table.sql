-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user"(
  id SERIAL PRIMARY KEY,
  
  user_name VARCHAR(64) UNIQUE NOT NULL,
  email VARCHAR (128) UNIQUE NOT NULL,
  phone VARCHAR (20) UNIQUE NOT NULL,
  password TEXT NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  deleted_at TIMESTAMPTZ DEFAULT NULL

);

--INDEXES HERE
CREATE INDEX idx_user__deleted_at ON "user" (deleted_at);
CREATE INDEX idx_user_unique_active_email on "user" (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_unique_active_user_name on "user" (user_name) WHERE deleted_at IS NULL;


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
