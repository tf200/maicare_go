
SEED_ORGANISATIONS ?= 6
SEED_LOCATIONS_PER_ORG ?= 2
SEED_SENDERS ?= 12
SEED_REGISTRATION_FORMS ?= 25
SEED_WAITING_LIST_CLIENTS ?= 12

migrateup:
	migrate -path db/migrations -database "postgresql://maicare:maicare@localhost:5432/maicare?sslmode=disable" -verbose up 1

migrateforce:
	migrate -path db/migrations -database "postgresql://maicare:maicare@localhost:5432/maicare?sslmode=disable" force 1


migratedown:
	migrate -path db/migrations -database "postgresql://maicare:maicare@localhost:5432/maicare?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockdb:
	mockgen -package mockdb -destination=db/mock/store.go github.com/rokunisan/chat_app/db/sqlc Store

swagger:
	swag init --parseDependency --output ./docs --generalInfo server.go --dir ./api

roles:
	go run cmd/roles/main.go

admin:
	go run cmd/admin/main.go

seed:
	go run cmd/seed/main.go -organisations $(SEED_ORGANISATIONS) -locations-per-org $(SEED_LOCATIONS_PER_ORG) -senders $(SEED_SENDERS) -count $(SEED_REGISTRATION_FORMS) -waiting-list-clients $(SEED_WAITING_LIST_CLIENTS)


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


.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mockdb swaggerrest_post roles admin seed push migrateforce update-proto generate-grpc lint mocks
