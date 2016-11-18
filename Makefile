#--------- RUN ON HOST -----------------
## Remove existing containers, build images and start db + user_api
stack-up:
	(docker stop postgres || exit 0) && \
	(docker rm postgres || exit 0) && \
	(docker stop user_api || exit 0) && \
  (docker rm user_api || exit 0) && \
	docker-compose build && \
	docker-compose up -d --force-recreate postgres && \
	echo 'Waiting for postgres to have started up..'; sleep 20 && \
	docker-compose up -d user_api

# Create images only
images:
	docker-compose build --force-rm

# Create DB container
create-db:
	(docker stop postgres || exit 0) && \
	(docker rm postgres || exit 0) && \
	docker-compose up -d --force-recreate postgres

# Create API container
create-api:
	(docker stop user_api || exit 0) && \
  (docker rm user_api || exit 0) && \
	docker-compose up -d --force-recreate user_api

## Tools
# Access pg shell
database-shell:
	docker exec -it postgres psql -Ucc cc_users

# Run DB migrations
migrate-db:
	docker exec -it user_api bash -c 'cd go/src/github.com/arbolista-dev/cc-user-api; goose -env  $(CC_ENV) up'

# Update API with latest changes (CC_ENV=dev only)
update-api:
	docker cp . user_api:/go/src/github.com/arbolista-dev/cc-user-api

# Backup database
backup-db:
	docker exec -it postgres pg_dump -c -U$(CC_DBUSER) $(CC_DBNAME) > ~/$(shell date  +'%Y%m%d-%H%M%S')-$(CC_DBNAME)-dump.sql

stack-test:
	docker-compose build && \
	docker-compose -f docker-compose.test.yml up -d postgres && \
	echo 'Waiting for postgres to have started up..'; sleep 30 && \
	docker-compose -f docker-compose.test.yml up user_api
