-- cat sql/0003-addCountryMigration.sql | PGPASSWORD=pass psql -h localhost -Ucc cc_users

ALTER TABLE IF EXISTS "users" ADD COLUMN "country" VARCHAR(2);
