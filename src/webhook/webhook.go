package webhook

import (
	"errors"
	"log"
	"os"
	"osoba/deploy"

	"golang.org/x/crypto/bcrypt"
)

type Info struct {
	*deploy.Info
	Token string
}

func FetchInfo(path string) (Info, error) {
	i := Info{
		Info: &deploy.Info{
			Path:       "/aaa",
			RootPath:   "/www/html",
			ReleaseURL: "https://github.com/w-haibara/portfolio/releases/download/v1.0.8/portfolio.zip",
		},
		Token: os.Getenv("OSOBA_TMP_TOKEN"),
	}

	if path == i.Path {
		return i, nil
	}

	return Info{}, errors.New("unknown error")
}

func (i Info) KeyVerify(authHeader []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(i.Token), authHeader); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
