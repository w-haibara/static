package auth

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	LoginFormURI  string
	VerifyKeyFile string
	CookieName    string
}

type App struct {
	*Config
	Claims    jwt.MapClaims
	verifyKey []byte
}

func NewApp(c Config) (App, error) {
	if c.LoginFormURI == "" {
		c.LoginFormURI = "/login?backTo=/osoba"
	}
	if c.VerifyKeyFile == "" {
		c.VerifyKeyFile = "/jwt-secret/secret.key"
	}
	if c.CookieName == "" {
		c.CookieName = "jwt_token"
	}

	a := App{
		Config:    &c,
		Claims:    make(jwt.MapClaims),
		verifyKey: []byte{},
	}

	var err error
	a.verifyKey, err = ioutil.ReadFile(a.VerifyKeyFile)
	if err != nil {
		return App{}, err
	}

	return a, nil
}

func (a App) Auth(next http.Handler, w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie(a.CookieName)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(c.Value, func(*jwt.Token) (interface{}, error) {
		return a.verifyKey, nil
	})
	if err != nil {
		return err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for k, v := range claims {
			a.Claims[k] = v
		}
		return nil
	}

	return errors.New("unknown error")
}
