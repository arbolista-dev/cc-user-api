
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP INDEX IF EXISTS "leaders_public_footprint_index";
DROP VIEW IF EXISTS "leaders_public_food_footprint";
DROP VIEW IF EXISTS "leaders_public_housing_footprint";
DROP VIEW IF EXISTS "leaders_public_shopping_footprint";
DROP VIEW IF EXISTS "leaders_public_transport_footprint";

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
