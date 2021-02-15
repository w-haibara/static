package main

import (
	"log"
	"net/http"
	"osoba/auth"
)

func main() {
	authConfig, err := auth.InitConfig(auth.Config{
		LoginFormURI:  "/login?backTo=/osoba",
		VerifyKeyFile: "/jwt-secret/secret.key",
	})
	if err != nil {
		log.Panic(err)
	}

	http.Handle("/", loggingHandler(checkMethodHandler(http.MethodGet, authHandler(authConfig, http.HandlerFunc(mainHandler)))))
	http.Handle("/deploy", loggingHandler(checkMethodHandler(http.MethodPost, webhookHandler())))

	http.ListenAndServe(":8080", nil)
}
