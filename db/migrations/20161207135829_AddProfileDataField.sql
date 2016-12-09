
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin

DO $$
    BEGIN
        BEGIN
            ALTER TABLE IF EXISTS "users" ADD COLUMN "profile_data" JSONB;
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column profile_data already exists in users.';
        END;
    END;
$$
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "profile_data";
