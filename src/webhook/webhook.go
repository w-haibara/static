package webhook

import (
	"log"
	"net/http"
	"os"
	"osoba/deploy"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	*deploy.Config
}

func InitConfigs(configs []Config) []Config {
	return configs
}

func (config Config) KeyVerify(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")

	if err := bcrypt.CompareHashAndPassword(getToken(config.Path), []byte(authHeader)); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func getToken(path string) []byte {
	// TODO: using key value store
	return []byte(os.Getenv("OSOBA_TOKEN_" + path))
}
