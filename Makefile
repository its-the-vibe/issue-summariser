.PHONY: build clean test install

# Build the application
build:
	go build -o issue-summariser main.go

# Install dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	rm -f issue-summariser

# Install the binary to $GOPATH/bin
install:
	go install

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Build and run
run: build
	@echo "Usage: echo '{\"message\": \"your message\"}' | ./issue-summariser"
