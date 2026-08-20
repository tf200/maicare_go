
INFISICAL_ENV ?= local
MIGRATION_VERSION ?= 1
LOCAL_DB_URL ?= postgresql://maicare:maicare@127.0.0.1:5432/maicare?sslmode=disable
REMOTE_DEV_DB_URL ?= postgresql://maicare:maicare@167.86.75.250:25432/maicare?sslmode=disable

SEED_FLAGS_LOCAL ?=
SEED_FLAGS_REMOTE ?= -organisations=1 -locations-per-org=1 -departments=1 -handbook-templates-per-department=1 -employee-handbook-assignments-per-department=0 -senders=3 -count=5 -other-intake-forms=0 -waiting-list-clients=2 -in-care-clients=1 -out-of-care-clients=0 -evaluations-per-in-care-client=1 -diagnoses-per-client=1 -medication-orders-per-client=1 -incidents-per-client=1 -invoices-per-client=0 -payments-per-invoice=0 -seed=42

db-local: migrate-up-local roles-sync-local seed-local admin-local
db-remote: migrate-up-remote roles-sync-remote seed-remote admin-remote

seed-local:
	DB_SOURCE="$(LOCAL_DB_URL)" go run ./cmd/seed $(SEED_FLAGS_LOCAL)

seed-remote:
	DB_SOURCE="$(REMOTE_DEV_DB_URL)" go run ./cmd/seed $(SEED_FLAGS_REMOTE)

migrate-up-local:
	migrate -path db/migrations -database "$(LOCAL_DB_URL)" -verbose up

migrate-up-remote:
	migrate -path db/migrations -database "$(REMOTE_DEV_DB_URL)" -verbose up

migrate-down-local:
	migrate -path db/migrations -database "$(LOCAL_DB_URL)" -verbose down 1

migrate-down-remote:
	migrate -path db/migrations -database "$(REMOTE_DEV_DB_URL)" -verbose down 1

migrate-force-local:
	migrate -path db/migrations -database "$(LOCAL_DB_URL)" force $(MIGRATION_VERSION)

migrate-force-remote:
	migrate -path db/migrations -database "$(REMOTE_DEV_DB_URL)" force $(MIGRATION_VERSION)

roles-sync-local:
	DB_SOURCE="$(LOCAL_DB_URL)" go run cmd/roles/main.go

roles-sync-remote:
	DB_SOURCE="$(REMOTE_DEV_DB_URL)" go run cmd/roles/main.go

admin-local:
	DB_SOURCE="$(LOCAL_DB_URL)" go run cmd/admin/main.go

admin-remote:
	DB_SOURCE="$(REMOTE_DEV_DB_URL)" go run cmd/admin/main.go

migrateup: migrate-up
migratedown: migrate-down
migrateforce: migrate-force
migrate-up: migrate-up-local
migrate-down: migrate-down-local
migrate-force: migrate-force-local
migrateup-local: migrate-up-local
migratedown-local: migrate-down-local
migrateforce-local: migrate-force-local
migrateup-remote: migrate-up-remote
migratedown-remote: migrate-down-remote
migrateforce-remote: migrate-force-remote
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

docker-up:
	@infisical export --env=$(INFISICAL_ENV) --format=dotenv > .env 2>/dev/null || \
	(echo "⚠️ Infisical unreachable / offline. Falling back to local app.env" && cp app.env .env)
	docker compose up -d --build

docker-down:
	docker compose down

env-sync:
	@echo "Syncing secrets from Infisical ($(INFISICAL_ENV)) into app.env..."
	infisical export --env=$(INFISICAL_ENV) --format=dotenv > app.env

server:
	@infisical run --env=$(INFISICAL_ENV) -- go run main.go 2>/dev/null || go run main.go

mockdb:
	mockgen -package mockdb -destination=db/mock/store.go github.com/rokunisan/chat_app/db/sqlc Store

swagger:
	swag init --parseDependency --output ./docs --generalInfo server.go --dir ./api

roles:
	$(MAKE) roles-sync-local

admin:
	$(MAKE) admin-local

seed:
	$(MAKE) seed-local


push:
	sudo docker build -t taha541/maicare:back . && sudo docker push taha541/maicare:back && git push

update-proto:
	git submodule update --remote --merge

generate-grpc:
	protoc \
		--go_out=grpclient --go_opt=paths=source_relative \
		--go-grpc_out=grpclient --go-grpc_opt=paths=source_relative \
		proto/service.proto proto/spelling_service.proto proto/reports_service.proto proto/schedule_service.proto

lint:
	golangci-lint run

mocks:
	go generate ./...


.PHONY: db-local db-remote seed-local seed-remote migrate-up-local migrate-up-remote migrate-down-local migrate-down-remote migrate-force-local migrate-force-remote roles-sync-local roles-sync-remote admin-local admin-remote migrate-up migrate-down migrate-force migrateup migratedown migrateforce migrateup-local migratedown-local migrateforce-local migrateup-remote migratedown-remote migrateforce-remote sqlc test server mockdb swagger roles admin seed push update-proto generate-grpc lint mocks docker-up docker-down env-sync
