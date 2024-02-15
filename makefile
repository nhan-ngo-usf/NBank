
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
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto
.PHONY: createdb dropdb migrateup migratedown migrateupuser migratedownuser sqlc postgresql test server mock proto