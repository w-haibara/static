osoba: *.go cmd/osoba/*.go cli/*.go
	gofmt -w *.go cmd/osoba/*.go cli/*.go
	go build ./cmd/...

.PHONY: init
init:
	go mod init osoba
	go mod tidy

.PHONY: test
test:
	gofmt -w *.go cmd/osoba/*.go cli/*.go
	go test -v ./...
