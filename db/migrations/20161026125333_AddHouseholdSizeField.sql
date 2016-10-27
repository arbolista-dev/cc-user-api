
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin

DO $$
    BEGIN
        BEGIN
            ALTER TABLE IF EXISTS "users" ADD COLUMN "household_size" SMALLINT;
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column household_size already exists in users.';
        END;
    END;
$$
-- +goose StatementEnd
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "household_size";
