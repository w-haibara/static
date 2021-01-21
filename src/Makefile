osoba: *.go go.mod
	gofmt -w *.go
	go build -o osoba main.go

.PHONY: init
init:
	go mod init osoba

.PHONY: test
test:
	gofmt -w *.go
	go test
