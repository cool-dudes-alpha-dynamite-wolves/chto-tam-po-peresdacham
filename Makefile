run:
	go run cmd/main.go

lint:
	golangci-lint run -v --config=.golangci.yaml
