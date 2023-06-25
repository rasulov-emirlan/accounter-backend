dev:
	go run cmd/apiserver/main.go --dev --migrate --env .env

migrate_new:
	goose -dir internal/storage/postgresql/migrations create $(name) sql