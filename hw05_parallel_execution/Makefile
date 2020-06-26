test: lint
	go test -v -race -count=100 .

lint: 
	golangci-lint run ./...
	