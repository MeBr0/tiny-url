.PHONY:
.SILENT:
.DEFAULT_GOAL := build

swag:
	go get github.com/swaggo/swag/cmd/swag
	swag init -g internal/app/app.go

lint:
	golangci-lint run

test:
	GIN_MODE=release go test --short -coverprofile=cover.out -v ./...

cover: test
	go tool cover -func=cover.out

build: swag
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app/main.go

run: build
	./app

fmt:
	gofmt -s -w .
