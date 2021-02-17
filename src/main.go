package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"osoba/config"

	"github.com/k0kubun/pp"
)

var configFile = flag.String("f", "config.json", "path to the configuration file")

func main() {
	flag.Parse()

	c := configure()

	http.Handle("/", loggingHandler(checkMethodHandler(http.MethodGet, authHandler(*c.Auth, http.HandlerFunc(mainHandler)))))
	http.Handle("/deploy", loggingHandler(checkMethodHandler(http.MethodPost, webhookHandler())))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func configure() config.Config {
	pp.Printf("config file: %v\n", *configFile)

	json, err := ioutil.ReadFile(*configFile)
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
