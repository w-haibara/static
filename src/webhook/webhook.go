package webhook

import (
	"errors"
	"log"
	"os"
	"osoba/deploy"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	*deploy.Config
	Token string
}

func FetchConfig(path string) (Config, error) {
	config := Config{
		Config: &deploy.Config{
			Path:       "/aaa",
			RootPath:   "/www/html",
			ReleaseURL: "https://github.com/w-haibara/portfolio/releases/download/v1.0.8/portfolio.zip",
		},
		Token: os.Getenv("OSOBA_TMP_TOKEN"),
	}

	if path == config.Path {
		return config, nil
	}

	return Config{}, errors.New("unknown error")
}

func (config Config) KeyVerify(authHeader []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(config.Token), authHeader); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
