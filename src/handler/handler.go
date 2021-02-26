package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"osoba/auth"
	"osoba/deploy"
	"osoba/resource"
	"osoba/webhook"
	"strconv"
)

func CheckMethods(next http.Handler, methods ...string) http.Handler {
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- logging handler ---")
		log.Printf("%#v\n", r)
		next.ServeHTTP(w, r)
	})
}

func Auth(config auth.Config, next http.Handler) http.Handler {
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

func Webhook(config resource.Config, chanDeployInfo chan deploy.Info) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- webhook handler ---")

		if r.Header.Get("Content-Type") != "application/json" {
			log.Println("[StatusBadRequest]", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			log.Println("[StatusInternalServerError]", http.StatusInternalServerError, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := make([]byte, length)
		length, err = r.Body.Read(body)
		if err != nil && err != io.EOF {
			log.Println("[StatusInternalServerError]", http.StatusInternalServerError, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var jsonBody map[string]string
		err = json.Unmarshal(body[:length], &jsonBody)
		if err != nil {
			log.Println("[StatusInternalServerError]", http.StatusInternalServerError, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		wh, err := webhook.FetchInfo(config, jsonBody["path"])
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
