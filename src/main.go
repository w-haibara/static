package main

import (
	"log"
	"net/http"
	"osoba/auth"
)

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- logging handler ---")
		log.Printf("%#v\n", r)
		next.ServeHTTP(w, r)
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- main handler ---")
	w.Write([]byte("OK"))
}

func main() {
	auth.Configure(auth.Config{
		LoginFormURI:  "/login?backTo=/osoba",
		VerifyKeyFile: "/jwt-secret/secret.key",
	})

	mainHandler := http.HandlerFunc(mainHandler)

	http.Handle("/", loggingHandler(auth.Handler(mainHandler)))
	http.ListenAndServe(":8080", nil)
}
