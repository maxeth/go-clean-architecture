PWD = $(shell pwd)
ACCTPATH = $(PWD)\account
PGCONTAINERNAME=backend-memrizr-course_postgres-db_1 # replace backend-memrizr-course with the name of your directory, or simply copy this container name from `docker-compose ps`
API_SERVICE_NAME=account

init:
	docker-compose up -d postgres-db && \
	make create-keypair ENV=test && \
	make createdb && \
	make migrateup N= && \
	docker-compose down

docker-compose-rebuild:
	docker compose build account && docker compose up

compose:
	docker-compose up

restart-compose:
	docker-compose down && docker-compose up

restart-service: #&& docker-compose build $(NAME) 
	docker-compose stop -t 1 $(NAME) && docker-compose up --no-start $(NAME) && docker-compose start $(NAME)

restart-api: #  && docker-compose build $(API_SERVICE_NAME) 
	docker-compose stop -t 1 $(API_SERVICE_NAME) && docker-compose up --no-start $(API_SERVICE_NAME) && docker-compose start $(API_SERVICE_NAME)

mock: 
	mockgen -package mocks -destination ./model/mocks/user_service.go github.com/maxeth/go-account-api/model UserRepository,UserService,TokenService,TokenRepository

gqlgen:
	go run github.com/99designs/gqlgen generate

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in rsa_private_$(ENV).pem -pubout -out rsa_public_$(ENV).pem 	

create-db:
	docker exec -u postgres $(PGCONTAINERNAME) createdb --username=postgres --owner=postgres accounts_db

drop-db:
	docker restart $(PGCONTAINERNAME) && docker exec -u postgres $(PGCONTAINERNAME) dropdb accounts_db

migrate-create:
	migrate create -ext sql -dir ./migration -seq -digits 5 postgres

migrate-up:
	migrate -path ./migration -database "postgresql://postgres:password@localhost:5432/accounts_db?sslmode=disable" -verbose up $(N)
 
migrate-down:
	migrate -path ./migration -database "postgresql://postgres:password@localhost:5432/accounts_db?sslmode=disable" -verbose down $(N)

