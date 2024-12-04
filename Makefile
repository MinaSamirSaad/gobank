# Build the Go binary
build:
	@go build -o bin/gobank

# Run the application in Docker and start the necessary services
run: build
	@./runDocker.sh
	@export JWTSecret=secret && ./bin/gobank

# Run tests for the Go application
test:
	@go test -v ./...
