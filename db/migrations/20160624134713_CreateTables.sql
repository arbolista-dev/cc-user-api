
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "users" (
    "user_id"          SERIAL PRIMARY KEY,
    "email"            VARCHAR(80) UNIQUE,
    "first_name"       VARCHAR(80),
    "last_name"        VARCHAR(80),
    "hash"             BYTEA,
    "salt"             BYTEA,
    "valid_jti"        BYTEA,
    "answers"          JSONB,
    "reset_hash"       BYTEA,
    "reset_expiration" TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "users";
