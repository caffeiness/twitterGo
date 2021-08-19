code-check: mod imports fmt vet
vet:
	go vet ./...
fmt:
	gofmt -d -s .
imports:
	goimports -w .
build:
	sh docker-build.sh
mod:
	go mod tidy
	go mod verify
	go mod download


