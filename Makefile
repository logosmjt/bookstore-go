DB_URL=postgresql://root:secret@localhost:5432/book-store?sslmode=disable

MIGRATION_PATH=db/migration

migrate_init:
	migrate create -ext sql -dir "$(MIGRATION_PATH)" -seq init_schema

migrate_up:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose down

.PHONY: migrate_init migrate_up migrate_down