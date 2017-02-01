
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin

DO $$
    BEGIN
        BEGIN
            ALTER TABLE IF EXISTS "users" ADD COLUMN "photo_url" VARCHAR(92);
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column photo_url already exists in users.';
        END;
    END;
$$
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "photo_url";
