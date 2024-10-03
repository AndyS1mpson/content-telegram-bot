LOCAL_MIGRATION_DIR=./migrations
POSTGRES = ${POSTGRES_URL}

# работа с базой данных
run-database:
	docker-compose -f ./deployments/docker-compose.yaml up --build

# Секция работы с миграциями
migrate:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) status && \
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) up
