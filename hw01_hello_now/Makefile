test: lint
	go test -v -count=1 -race -gcflags=-l -timeout=30s .

lint: 
	golangci-lint run ./...

build: test
	go build -o hw1 