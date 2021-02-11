osoba: *.go */*.go go.mod
	gofmt -w *.go */*.go
	go build -o osoba main.go handler.go

.PHONY: init
init:
	go mod init osoba

.PHONY: test
test:
	gofmt -w *.go
	go test
