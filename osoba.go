package osoba

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	DocumentRoot         string
	TmpDirContentsPrefix string
	Contents             Contents
}

type Content struct {
	URL    string
	Secret string
}

type Contents struct {
	Mu sync.RWMutex
	V  map[string]Content
}

func (a App) DeployHandler(chanDeployPath chan<- string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		log.Println("calling deploy handler, path:", path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("[StatusUnauthorized]", http.StatusUnauthorized, "Authorization header is empty")
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		if err := a.KeyVerify(path, []byte(authHeader)); err != nil {
			log.Println("[StatusUnauthorized]", http.StatusUnauthorized, "API key verify error:", err.Error())
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		chanDeployPath <- path

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Accepted\n"))
	})
}

func (a App) KeyVerify(path string, key []byte) error {
	if _, ok := a.Contents.V[path]; !ok {
		return fmt.Errorf("content not exist: " + path)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Contents.V[path].Secret), key); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
