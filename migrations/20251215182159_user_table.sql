-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_status_choise AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE user_role_choise AS ENUM ('admin', 'user', 'moderator');

CREATE TABLE IF NOT EXISTS "user"(
  id SERIAL PRIMARY KEY,
  uuid UUID UNIQUE NOT NULL,

  user_name VARCHAR(64) UNIQUE NOT NULL,
  email VARCHAR (128) UNIQUE NOT NULL,
  phone VARCHAR (20) UNIQUE NOT NULL,
  password TEXT NOT NULL,
  user_status user_status_choise NOT NULL DEFAULT 'active',
  user_role user_role_choise NOT NULL DEFAULT 'user',

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
DROP TYPE IF EXISTS user_status_choise;
DROP TYPE IF EXISTS user_role_choise;
-- +goose StatementEnd
