package webhook

import (
	"errors"
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
	if !cmpPwd(getToken(config.Path), []byte(authHeader)) {
		return errors.New("unknown error")
	}
	
	return nil
}

func getToken(path string) []byte {
	// TODO: using key value store
	return []byte(os.Getenv("OSOBA_TOKEN_" + path))
}

func cmpPwd(hashedPwd []byte, pwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
