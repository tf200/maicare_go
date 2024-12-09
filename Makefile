
migrateup:
	migrate -path db/migrations -database "postgresql://user:jCjVBejKHMCUC5VjtLbTz69hVgJBfUHC@dpg-ct8ncs68ii6s73cd8oug-a.oregon-postgres.render.com/mdev_z2v5" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://user:jCjVBejKHMCUC5VjtLbTz69hVgJBfUHC@dpg-ct8ncs68ii6s73cd8oug-a.oregon-postgres.render.com/mdev_z2v5" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockdb:
	mockgen -package mockdb -destination=db/mock/store.go github.com/rokunisan/chat_app/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mockdb rest_post