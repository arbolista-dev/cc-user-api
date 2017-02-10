
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE user_goals ALTER COLUMN status TYPE VARCHAR(40);
DROP TYPE IF EXISTS "status_options";

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
CREATE TYPE status_options AS ENUM ('pledged', 'completed', 'not_relevant');
ALTER TABLE user_goals ALTER COLUMN status TYPE status_options;
