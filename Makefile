build:
	@go build -o bin/ecommerce_backend_go cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecommerce_backend_go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down