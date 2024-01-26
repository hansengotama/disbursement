include .env

migrate-up:
	migrate -path ./sql/schema -database "postgres://${POSTGRES_DB_USER}:${POSTGRES_DB_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB_NAME}?sslmode=disable" up

migrate-rollback:
	migrate -path ./sql/schema -database "postgres://${POSTGRES_DB_USER}:${POSTGRES_DB_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB_NAME}?sslmode=disable" down

