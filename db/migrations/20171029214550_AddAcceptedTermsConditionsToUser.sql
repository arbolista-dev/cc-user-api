-- +goose Up
ALTER TABLE users ADD COLUMN accepted_terms_conditions TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN over_eighteen_years TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE users DROP COLUMN accepted_terms_conditions;
ALTER TABLE users DROP COLUMN over_eighteen_years;
