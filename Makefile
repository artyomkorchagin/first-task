include .env
export

db-up:
	@GOOSE_DRIVER=${GOOSE_DRIVER} 
	@GOOSE_DBSTRING="host=${SERVER_HOST} port=${SERVER_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSLMODE}" 
	@goose up

db-down:
	@GOOSE_DRIVER=${GOOSE_DRIVER} 
	@GOOSE_DBSTRING="host=${SERVER_HOST} port=${SERVER_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSLMODE}" 
	@goose down

dev-up:
	@docker-compose -f ./docker-compose/dev/docker-compose.yml up -d 

dev-down:
	@docker-compose -f ./docker-compose/dev/docker-compose.yml down -v

create-topic:
	@docker exec -it kafka kafka-topics.sh  --create  --topic orders --bootstrap-server localhost:9092 --partitions 1  --replication-factor 1