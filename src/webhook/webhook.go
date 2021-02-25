package webhook

import (
	"context"
	"fmt"
	"log"
	"osoba/db"
	"osoba/deploy"
	"osoba/resource"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const mongodbURI string = "mongodb://mongo:27017"

type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Path       string             `bson:"path,omitempty"`
	RootPath   string             `bson:"rootpath,omitempty"`
	ReleaseURL string             `bson:"releaseurl,omitempty"`
	Token      string             `bson:"token,omitempty"`
}

type Info struct {
	*deploy.Info
	Token string
}

func Init() {
	log.Println("DB initializing")
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, mongodbURI)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Panic(err)
		}
	}()

	pathes := client.NewDB("osoba").NewCollection("pathes", Document{})

	log.Println("delete collection")
	if err := pathes.Delete(ctx, bson.M{}); err != nil {
		log.Panic(err)
	}

	log.Println("set tmp data")
	// set tmp data
	if err := pathes.Insert(ctx, []Document{
		Document{
			Path:       "/aaa",
			RootPath:   "/www/html",
			ReleaseURL: "https://github.com/w-haibara/portfolio/releases/download/v1.0.8/portfolio.zip",
			Token:      "$2a$10$sIKCSbHCLnNALUnaeMg1muyPAb4wrM57xJ1sHYmuHhoUtz0u9cqR2",
		},
	}); err != nil {
		log.Panic(err)
	}

	var docs []Document
	if err := pathes.Read(ctx, bson.M{}, &docs); err != nil {
		log.Panic(err)
	}

	log.Println("DB initializing success")
}

func FetchInfo(c resource.Config, path string) (Info, error) {
	docs, err := c.Fetch(bson.M{"path": path})
	if err != nil {
		return Info{}, err
	}

	if len(docs) != 1 {
		return Info{}, fmt.Errorf("docs length is invalid: '%#v'\n", docs)
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
