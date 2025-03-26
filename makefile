install:
	go mod tidy

test:
	go clean -testcache
	go test ./...

test-v:
	go test -v ./...  # Modo verbose

test-cover:
	go clean -testcache
	go test -short -coverprofile=coverage.out ./... 2>&1
	go tool cover -func=coverage.out