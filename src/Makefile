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
	go-callvis -format=png -focus main -nostd -file=../graph/main ./ | dot -Tpng -o ../graph/main.png
	go-callvis -format=png -focus osoba/config -nostd -file=../graph/config ./ | dot -Tpng -o ../graph/config.png
	go-callvis -format=png -focus osoba/handler -nostd -file=../graph/handler ./ | dot -Tpng -o ../graph/handler.png
	go-callvis -format=png -focus osoba/auth -nostd -file=../graph/auth ./ | dot -Tpng -o ../graph/auth.png
	go-callvis -format=png -focus osoba/webhook -nostd -file=../graph/webhook ./ | dot -Tpng -o ../graph/webhook.png
	go-callvis -format=png -focus osoba/deploy -nostd -file=../graph/deploy ./ | dot -Tpng -o ../graph/deploy.png
	go-callvis -format=png -focus osoba/resource -nostd -file=../graph/resource ./ | dot -Tpng -o ../graph/resource.png
	go-callvis -format=png -focus osoba/db -nostd -file=../graph/db ./ | dot -Tpng -o ../graph/db.png
