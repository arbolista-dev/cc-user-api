#--------- RUN ON HOST -----------------
## Initial setup commands
# Create images for containers
images:
	docker pull postgres:9.4.5 && \
	docker build -t cc:user_api .

# Create DB container
create-database:
	(docker stop postgres || exit 0) && \
	(docker rm postgres || exit 0) && \
	docker run \
		-d \
		-p 127.0.0.1:15432:5432 \
		--name postgres postgres:9.4.5 && \
	sleep 15 && \
	docker exec postgres psql -h127.0.0.1 -p5432 -Upostgres -c "CREATE ROLE $(CC_DBUSER) PASSWORD '$(CC_DBPASS)' NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT LOGIN" &&\
	docker exec postgres psql -h127.0.0.1 -p5432 -Upostgres -c "CREATE DATABASE $(CC_DBNAME)" &&\
	cat sql/0001-init.sql sql/0002-leadersMigration.sql sql/0003-addCountryMigration.sql sql/0004-facebookMigration.sql | PGPASSWORD=$(CC_DBPASS) psql -h127.0.0.1 -p15432 -U$(CC_DBUSER) $(CC_DBNAME)

# Create API container
create-api:
	(docker stop user_api || exit 0) && \
  (docker rm user_api || exit 0) && \
	docker run \
		-d \
		-p 0.0.0.0:8082:8082 \
		--name user_api\
		--link postgres \
		--env-file .env\
		cc:user_api

## Tools
# Access pg shell
database-shell:
	docker exec -it postgres psql -Ucc cc_users

# Clear database and run migrations
reset-database:
	cat sql/0000-reset.sql sql/0001-init.sql sql/0002-leadersMigration.sql sql/0003-addCountryMigration.sql sql/0004-facebookMigration.sql | PGPASSWORD=$(CC_DBPASS) psql -h127.0.0.1 -p15432 -U$(CC_DBUSER) $(CC_DBNAME)

# Run DB migrations
migrate-database:
	cat sql/0001-init.sql sql/0002-leadersMigration.sql sql/0003-addCountryMigration.sql sql/0004-facebookMigration.sql | PGPASSWORD=$(CC_DBPASS) psql -h127.0.0.1 -p15432 -U$(CC_DBUSER) $(CC_DBNAME)

# Update API with latest changes
update-api:
	docker cp . user_api:/go/src/github.com/arbolista-dev/cc-user-api
