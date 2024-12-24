
createdb:
	docker exec -it postgresDB createdb --username=myuser simple_bank
dropdb:
	docker exec -it postgresDB dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateupuser:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedownuser:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	docker-compose run --rm sqlc
postgresql:
	docker-compose up -d postgresDB
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination=db/mock/store.go github.com/nhan-ngo-usf/NBank/db/sqlc Store
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank \
	proto/*.proto
	statik -src=./doc/swagger -dest ./doc
evans:
	evans -r repl --host localhost --port=8080
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
.PHONY: createdb dropdb migrateup migratedown migrateupuser migratedownuser sqlc postgresql test server mock proto evans redis