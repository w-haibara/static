package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"osoba/config"

	"github.com/k0kubun/pp"
)

func main() {
	c := configure()

	http.Handle("/", loggingHandler(checkMethodHandler(http.MethodGet, authHandler(*c.Auth, http.HandlerFunc(mainHandler)))))
	http.Handle("/deploy", loggingHandler(checkMethodHandler(http.MethodPost, webhookHandler())))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func configure() config.Config {
	json, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panic(err)
	}

	c, err := config.Configure(json)
	if err != nil {
		log.Panic(err)
	}

	pp.Println("Auth:", c.Auth)

	return c
}
