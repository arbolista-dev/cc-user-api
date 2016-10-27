
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP VIEW IF EXISTS "leaders_public_footprint";
CREATE VIEW leaders_public_footprint AS
  SELECT first_name, last_name, total_footprint, city, state, county, household_size
  FROM users
  WHERE public IS TRUE
    AND (total_footprint->'result_grand_total') IS NOT NULL
    AND household_size IS NOT NULL
  ORDER BY total_footprint->'result_grand_total' ASC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP VIEW IF EXISTS "leaders_public_footprint";
