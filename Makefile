run:
	go run cmd/main.go
test:
	go test ./tests/... -v
test-integration:
	TEST_TYPE=integration go test ./tests/integration/... -v


swagger:
	swag fmt && swag init -g ./cmd/main.go -o ./docs --parseInternal=true --parseDependency=true
