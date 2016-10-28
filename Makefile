#--------- RUN ON HOST -----------------
## Initial setup commands
# Create images for containers
images:
	docker pull postgres:9.4.5 && \
	docker build -t cc:user_api .
	go get -u bitbucket.org/liamstask/goose/cmd/goose

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
	goose up

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

# Run DB migrations
migrate-database:
	goose up

# Update API with latest changes
update-api:
	docker cp . user_api:/go/src/github.com/arbolista-dev/cc-user-api
