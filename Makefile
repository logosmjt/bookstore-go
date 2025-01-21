DB_URL=postgresql://root:secret@localhost:5432/book-store?sslmode=disable

MIGRATION_PATH=db/migration

migrate_init:
	migrate create -ext sql -dir "$(MIGRATION_PATH)" -seq init_schema

migrate_up:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/logosmjt/bookstore-go/db/sqlc Store

.PHONY: migrate_init migrate_up migrate_down sqlc server mock test