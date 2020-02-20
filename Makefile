test:
	go test -race -v ./...

lint:
	@test -z $(gofmt -s -l -d ./)
	@golangci-lint run ./
