
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
DO $$
    BEGIN
        BEGIN
            CREATE VIEW leaders_public_footprint AS
              SELECT first_name, last_name, total_footprint->'result_grand_total' as footprint, city, state, county
              FROM users
              WHERE public IS TRUE AND (total_footprint->'result_grand_total') IS NOT NULL
              ORDER BY total_footprint->'result_grand_total' ASC;
        EXCEPTION
            WHEN duplicate_table THEN RAISE NOTICE 'view leaders_public_footprint already exists.';
        END;

        BEGIN
            CREATE VIEW leaders_public_food_footprint AS
              SELECT first_name, last_name, total_footprint->'result_food_total' as footprint, city, state, county
              FROM users
              WHERE public IS TRUE AND (total_footprint->'result_food_total') IS NOT NULL
              ORDER BY total_footprint->'result_food_total' ASC;
        EXCEPTION
            WHEN duplicate_table THEN RAISE NOTICE 'view leaders_public_food_footprint already exists.';
        END;

        BEGIN
            CREATE VIEW leaders_public_housing_footprint AS
              SELECT first_name, last_name, total_footprint->'result_housing_total' as footprint, city, state, county
              FROM users
              WHERE public IS TRUE AND (total_footprint->'result_housing_total') IS NOT NULL
              ORDER BY total_footprint->'result_housing_total' ASC;
        EXCEPTION
            WHEN duplicate_table THEN RAISE NOTICE 'view leaders_public_housing_footprint already exists.';
        END;

        BEGIN
            CREATE VIEW leaders_public_shopping_footprint AS
              SELECT first_name, last_name, (total_footprint->'result_shopping_total') as footprint, city, state, county
              FROM users
              WHERE public IS TRUE AND (total_footprint->'result_shopping_total') IS NOT NULL
              ORDER BY total_footprint->'result_shopping_total' ASC;
        EXCEPTION
            WHEN duplicate_table THEN RAISE NOTICE 'view leaders_public_shopping_footprint already exists.';
        END;

        BEGIN
            CREATE VIEW leaders_public_transport_footprint AS
              SELECT first_name, last_name, (total_footprint->'result_transport_total') as footprint, city, state, county
              FROM users
              WHERE public IS TRUE AND (total_footprint->'result_transport_total') IS NOT NULL
              ORDER BY total_footprint->'result_transport_total' ASC;
        EXCEPTION
            WHEN duplicate_table THEN RAISE NOTICE 'view leaders_public_transport_footprint already exists.';
        END;
    END;
$$
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP VIEW IF EXISTS "leaders_public_footprint";
DROP VIEW IF EXISTS "leaders_public_food_footprint";
DROP VIEW IF EXISTS "leaders_public_housing_footprint";
DROP VIEW IF EXISTS "leaders_public_shopping_footprint";
DROP VIEW IF EXISTS "leaders_public_transport_footprint";
