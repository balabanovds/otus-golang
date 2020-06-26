build: test clean
	go build -o build/go-dd

test: lint
	go test -v -count=1 -race -gcflags=-l -timeout=30s .

lint: 
	golangci-lint run ./...

clean:
	rm -rf build/