
migrateup:
	migrate -path db/migrations -database "postgresql://maicaredb_user:HxGiWC1fiFRIZZ8dxeTRpOGJUdyRIox4@dpg-cu56483tq21c73e1bqmg-a.frankfurt-postgres.render.com/maicaredb" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://maicaredb_user:HxGiWC1fiFRIZZ8dxeTRpOGJUdyRIox4@dpg-cu56483tq21c73e1bqmg-a.frankfurt-postgres.render.com/maicaredb" -verbose down

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

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mockdb swaggerrest_post roles admin push