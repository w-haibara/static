osoba: *.go */*.go go.mod
	gofmt -w *.go */*.go
	go build -o osoba main.go

.PHONY: init
init:
	go mod init osoba

.PHONY: test
test:
	gofmt -w *.go
	go test

.PHONY: graph
graph:
	go-callvis -format=png -nostd -file=./graph/osoba ./ | dot -Tpng -o ./graph/osoba.png
