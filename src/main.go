package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"osoba/config"
	"osoba/deploy"
	"osoba/webhook"

	"github.com/k0kubun/pp"
)

var configFile = flag.String("f", "config.json", "path to the configuration file")

func main() {
	flag.Parse()

	c := configure()
	chanDeployInfo := make(chan deploy.Info)

	http.Handle("/", loggingHandler(checkMethodsHandler(authHandler(*c.Auth, http.HandlerFunc(mainHandler)), http.MethodGet)))
	http.Handle("/deploy", loggingHandler(checkMethodsHandler(webhookHandler(*c.DB, chanDeployInfo), http.MethodPost)))

	go deploy.AwaitDeploy(chanDeployInfo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func configure() config.Config {
	log.Printf("config file: %v\n", *configFile)

	json, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic(err)
	}

	c, err := config.Configure(json)
	if err != nil {
		log.Panic(err)
	}

	pp.Println(c)

	webhook.Init()

	return c
}
