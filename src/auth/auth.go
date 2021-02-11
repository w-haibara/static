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
	Claims        jwt.MapClaims

	verifyKey []byte
}

func InitConfig(c Config) (Config, error) {
	if c.CookieName == "" {
		c.CookieName = "jwt_token"
	}
	config := Config{
		LoginFormURI:  c.LoginFormURI,
		VerifyKeyFile: c.VerifyKeyFile,
		CookieName:    c.CookieName,
		Claims:        make(jwt.MapClaims),
	}

	var err error
	config.verifyKey, err = ioutil.ReadFile(config.VerifyKeyFile)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config Config) Auth(next http.Handler, w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie(config.CookieName)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(c.Value, func(*jwt.Token) (interface{}, error) {
		return config.verifyKey, nil
	})
	if err != nil {
		return err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for k, v := range claims {
			config.Claims[k] = v
		}
		return nil
	}

	return errors.New("unknown error")
}
