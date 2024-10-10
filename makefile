
createdb:
	docker exec -it postgresDB createdb --username=myuser simple_bank
dropdb:
	docker exec -it postgresDB dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	docker-compose run --rm sqlc
postgresql:
	docker-compose up -d postgresDB

.PHONY: createdb dropdb migrateup migratedown sqlc postgresql