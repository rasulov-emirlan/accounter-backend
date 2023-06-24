dev:
	go run cmd/apiserver/main.go --dev --migrate

migrate_new:
	goose -dir internal/storage/postgresql/migrations create $(name) sql