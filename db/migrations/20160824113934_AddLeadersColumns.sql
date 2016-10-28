
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin

DO $$
    BEGIN
        BEGIN
            ALTER TABLE IF EXISTS "users" ADD COLUMN "public" BOOLEAN,
                                          ADD COLUMN "city" VARCHAR(80),
                                          ADD COLUMN "state" VARCHAR(80),
                                          ADD COLUMN "county" VARCHAR(80),
                                          ADD COLUMN "total_footprint" JSONB;
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column public already exists in users.';
        END;
    END;
$$
-- +goose StatementEnd


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "public",
                              DROP COLUMN IF EXISTS "city",
                              DROP COLUMN IF EXISTS "state",
                              DROP COLUMN IF EXISTS "county",
                              DROP COLUMN IF EXISTS "total_footprint";
