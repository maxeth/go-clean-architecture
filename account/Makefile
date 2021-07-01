PWD = $(shell pwd)
ACCTPATH = $(PWD)\account
PGCONTAINERNAME=backend-memrizr-course_postgres-db_1 # replace backend-memrizr-course with the name of your directory, or simply copy this container name from `docker-compose ps`

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

restart:
	docker-compose down && docker-compose up

mock: 
	mockgen -package mocks -destination ./model/mocks/user_service.go github.com/maxeth/go-account-api/model UserRepository,UserService,TokenService,TokenRepository


create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in rsa_private_$(ENV).pem -pubout -out rsa_public_$(ENV).pem 	

createdb:
	docker exec -u postgres $(PGCONTAINERNAME) createdb --username=postgres --owner=postgres accounts_db

dropdb:
	docker restart $(PGCONTAINERNAME) && docker exec -u postgres $(PGCONTAINERNAME) dropdb accounts_db

migratecreate:
	migrate create -ext sql -dir ./migration -seq -digits 5 postgres

migrateup:
	migrate -path ./migration -database "postgresql://postgres:password@localhost:5432/accounts_db?sslmode=disable" -verbose up $(N)
 
migratedown:
	migrate -path ./migration -database "postgresql://postgres:password@localhost:5432/accounts_db?sslmode=disable" -verbose down $(N)

