package webhook

import (
	"context"
	"fmt"
	"log"
	"os"
	"osoba/db"
	"osoba/deploy"
	"time"

	"github.com/k0kubun/pp"
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

	var docs []Document
	if err := pathes.Read(ctx, bson.M{"path": "/aaa"}, &docs); err != nil {
		log.Panic(err)
	}
	if len(docs) == 1 {
		log.Println("DB status OK")
		return
	}

	log.Println("insert data")
	if err := pathes.Insert(ctx, []Document{
		Document{
			Path:       "/aaa",
			RootPath:   "/www/html",
			ReleaseURL: "https://github.com/w-haibara/portfolio/releases/download/v1.0.8/portfolio.zip",
			Token:      os.Getenv("OSOBA_TMP_TOKEN"),
		},
	}); err != nil {
		log.Panic(err)
	}

	if err := pathes.Read(ctx, bson.M{}, &docs); err != nil {
		log.Panic(err)
	}
	pp.Println(docs)
	log.Println("DB initializing success")
}

func FetchInfo(path string) (Info, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, mongodbURI)
	if err != nil {
		return Info{}, err
	}
	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	pathes := client.NewDB("osoba").NewCollection("pathes", Document{})
	var docs []Document
	if err := pathes.Read(ctx, bson.M{"path": path}, &docs); err != nil {
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
