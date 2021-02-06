package auth

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var (
	config    Config = Config{}
	verifyKey []byte
)

type Config struct {
	LoginFormURI  string
	VerifyKeyFile string
}

func Configure(c Config) error {
	config.LoginFormURI = c.LoginFormURI
	config.VerifyKeyFile = c.VerifyKeyFile

	var err error
	verifyKey, err = ioutil.ReadFile(config.VerifyKeyFile)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- auth handler ---")

		if c, err := r.Cookie("jwt_token"); err == nil {
			token, err := jwt.Parse(c.Value, func(*jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("JWT parse error:", err.Error())
				return
			} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				log.Printf("%#v\n", claims)
				next.ServeHTTP(w, r)
				return
			}
		}

		log.Println("redirect to login form (", config.LoginFormURI, ")")
		http.Redirect(w, r, config.LoginFormURI, http.StatusSeeOther)
	})
}
