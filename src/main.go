package main

import (
	"log"
	"net/http"
	"osoba/auth"
	"osoba/logging"
)

func main() {
	auth.Configure(auth.Config{
		LoginFormURI:  "/login?backTo=/osoba",
		VerifyKeyFile: "/jwt-secret/secret.key",
	})

	mainHandler := http.HandlerFunc(mainHandler)

	http.Handle("/", logging.Handler(auth.Handler(mainHandler)))
	http.ListenAndServe(":8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- main handler ---")
	w.Write([]byte("OK"))
}
