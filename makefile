dev:
	go run main.go
build:
	go build
start: build
	./go-auth
test:
	go test ./...