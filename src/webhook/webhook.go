package webhook

import (
	"fmt"
	"log"
	"osoba/deploy"
	"osoba/resource"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const mongodbURI string = "mongodb://mongo:27017"

type Info struct {
	*deploy.Info
	Token string
}

func FetchInfo(c resource.Config, path string) (Info, error) {
	docs, err := c.Fetch(bson.M{"path": path})
	if err != nil {
		return Info{}, err
	}

	if len(docs) != 1 {
		return Info{}, fmt.Errorf("docs length is invalid, path: %s, docs: %#v", path, docs)
	}

	doc := docs[0]

	i := Info{
		Info: &deploy.Info{
			Path:       doc.Path,
			RootPath:   doc.RootPath,
			ReleaseURL: doc.ReleaseURL,
		},
		Token: doc.Token,
	}

	return i, nil
}

func (i Info) KeyVerify(authHeader []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(i.Token), authHeader); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
