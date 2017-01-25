
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TYPE status_options AS ENUM ('pledged', 'completed', 'not_relevant');

CREATE TABLE IF NOT EXISTS "actions" (
    "action_id"       SERIAL PRIMARY KEY,
    "key"             VARCHAR(80),
    "user_id"         INTEGER REFERENCES users(user_id),
    "status"          status_options,
    "created_at"      TIMESTAMP,
    "tons_saved"      REAL,
    "dollars_saved"   REAL,
    "upfront_cost"    REAL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "actions";
DROP TYPE IF EXISTS "status_options";
