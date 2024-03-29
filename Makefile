# DSN=host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5
BINARY_NAME=subscriptionService
REDIS="127.0.0.1:6379"

## build: builds all binaries
########################
build_run:
	@echo Building broker binary...
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go env -w CGO_ENABLED=0 && go build -o ${BINARY_NAME} ./cmd/web
	docker-compose up -d
	@echo Done!


test:
	@echo "Testing..."
	go test -v ./...



