package main

import (
	"log"
	"net/http"
	"osoba/auth"
	"osoba/webhook"
)

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- logging handler ---")
		log.Printf("%#v\n", r)
		next.ServeHTTP(w, r)
	})
}

func authHandler(config auth.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- auth handler ---")
		err := config.Auth(next, w, r)
		if err != nil {
			log.Println("JWT parse error:", err.Error())
			log.Println("redirect to login form (", config.LoginFormURI, ")")
			http.Redirect(w, r, config.LoginFormURI, http.StatusSeeOther)
			return
		}
		log.Printf("%#v\n", config.Claims)
		next.ServeHTTP(w, r)
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- main handler ---")
	w.Write([]byte("OK"))
}

func webhookHandler(config webhook.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- webhook handler ---")
		err := config.KeyVerify(w, r)
		if err != nil {
			log.Println("API key verify error:", err.Error())
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		if err := config.Deploy(); err != nil {
			log.Println("Deploy error:", err.Error())
			http.Error(w, "Deploy failed.", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Deploy succsess\n"))
	})

}

func webhooksManage(webhooks []webhook.Config) {
	for _, v := range webhooks {
		http.Handle(v.Path, loggingHandler(webhookHandler(v)))
	}
}
