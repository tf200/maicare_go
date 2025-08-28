
migrateup:
	migrate -path db/migrations -database "postgresql://maicaredb_user:HxGiWC1fiFRIZZ8dxeTRpOGJUdyRIox4@167.86.75.250:5432/maicaredb?sslmode=disable" -verbose up

migrateforce:
	migrate -path db/migrations -database "postgresql://maicaredb_user:HxGiWC1fiFRIZZ8dxeTRpOGJUdyRIox4@167.86.75.250:5432/maicaredb?sslmode=disable" force 1


migratedown:
	migrate -path db/migrations -database "postgresql://maicaredb_user:HxGiWC1fiFRIZZ8dxeTRpOGJUdyRIox4@167.86.75.250:5432/maicaredb?sslmode=disable" -verbose down

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
	cd roles && g++ -std=c++17 -o rbac_sync rbac_sync.cpp -lpqxx -lpq -lyaml-cpp && ./rbac_sync && cd ..
admin:
	cd admin && g++ -o admin admin.cpp -lpqxx -lssl -lcrypto -l:bcrypt.a && ./admin && cd ..


push:
	sudo docker build -t taha541/maicare:back . && sudo docker push taha541/maicare:back && git push

update-proto:
	git submodule update --remote --merge

generate-grpc:
	protoc \
		--go_out=grpclient --go_opt=paths=source_relative \
		--go-grpc_out=grpclient --go-grpc_opt=paths=source_relative \
		proto/service.proto proto/spelling_service.proto

lint:
	golangci-lint run

mocks:
	go generate ./...


.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mockdb swaggerrest_post roles admin push migrateforce update-proto generate-grpc lint mocks