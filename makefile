dev:
	go run cmd/main.go
build:
	go build cmd/main.go
start: build
	./main
test:
	go test ./...