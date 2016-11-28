
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
DO $$
    BEGIN
        BEGIN
            ALTER TABLE IF EXISTS "users" ADD COLUMN "email_hash"       BYTEA,
    									  ADD COLUMN "email_expiration" TIMESTAMP;
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column public already exists in users.';
        END;
    END;
$$
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "email_hash",
                              DROP COLUMN IF EXISTS "email_expiration";