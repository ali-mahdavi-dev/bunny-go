run:
	go run cmd/main.go
test:
	go test ./tests/... -v
test-integration:
	TEST_TYPE=integration go test ./tests/integration/... -v

run:
	go run cmd/main.go

migrate-up:
	go run ./cmd/main.go migrate up

migrate-down:
	go run ./cmd/main.go migrate down --env-file=./.env

swagger:
	swag fmt && swag init -g ./cmd/main.go -o ./docs --parseInternal=true --parseDependency=true
