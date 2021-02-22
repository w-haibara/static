package main

import (
	"fmt"
	"log"
	"net/http"
	"osoba/auth"
	"osoba/deploy"
	"osoba/webhook"
)

func checkMethodsHandler(next http.Handler, methods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- check methods handler ---")

		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}
		log.Println("[StatusMethodNotAllowed]", http.StatusMethodNotAllowed, "must:", methods, ", have:", r.Method)
		http.Error(w, "only:"+fmt.Sprint(methods), http.StatusMethodNotAllowed)
	})
}

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

		app, err := auth.NewApp(config)
		if err != nil {
			return
		}

		c, err := r.Cookie(app.CookieName)
		if err != nil {
			log.Println("read cookie error", err)
			return
		}

		if err := app.Auth(c.Value); err != nil {
			log.Println("JWT parse error:", err.Error())
			log.Println("redirect to login form (", app.LoginFormURI, ")")
			http.Redirect(w, r, app.LoginFormURI, http.StatusSeeOther)
			return
		}
		log.Printf("%#v\n", app.Claims)
		next.ServeHTTP(w, r)
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- main handler ---")
	w.Write([]byte("OK"))
}

func webhookHandler(chanDeployInfo chan deploy.Info) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- webhook handler ---")

		if err := r.ParseForm(); err != nil {
			log.Println("[StatusBadRequest]", http.StatusBadRequest, err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		wh, err := webhook.FetchInfo(r.PostFormValue("path"))
		if err != nil {
			log.Println("[Internal Server Error]", http.StatusInternalServerError, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("[StatusUnauthorized]", http.StatusUnauthorized, "Authorization header is empty")
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		if err := wh.KeyVerify([]byte(authHeader)); err != nil {
			log.Println("[StatusUnauthorized]", http.StatusUnauthorized, "API key verify error:", err.Error())
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		chanDeployInfo <- *wh.Info

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Accepted\n"))
	})

}
