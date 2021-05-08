package cli

import (
	"log"
	"net/http"
	"osoba"

	"github.com/k0kubun/pp"
)

func Run() int {
	log.Println("setting app config")
	a := osoba.App{}
	if err := a.LoadConfig(); err != nil {
		log.Fatal("[Failed to load config]", err)
	}
	pp.Println(a)

	chanDeployPath := make(chan string)

	http.Handle("/api/deploy", a.DeployHandler(chanDeployPath))

	go func() {
		for {
			path := <-chanDeployPath
			log.Println("deploy starting:", path)
			if err := a.Deploy(path); err != nil {
				log.Println(err)
				continue
			}
			log.Println("deploy complete:", path)
		}
	}()

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))

	return 255
}
