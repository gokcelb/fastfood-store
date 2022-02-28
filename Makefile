test:
	go test ./...

run:
	go run main.go

lint:
	golangci-lint run