DB_URL=postgresql://root:secret@localhost:5432/book-store?sslmode=disable

MIGRATION_PATH=db/migration

migrate_init:
	migrate create -ext sql -dir "$(MIGRATION_PATH)" -seq init_schema

migrate_up:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose down

migration_new:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/logosmjt/bookstore-go/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bookstore \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: migrate_init migrate_up migrate_down sqlc server mock test proto evans