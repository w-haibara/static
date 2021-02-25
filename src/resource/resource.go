package resource

import (
	"context"
	"log"
	"osoba/db"
	"time"

	"github.com/k0kubun/pp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	MongodbURI string
}

type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Path       string             `bson:"path,omitempty"`
	RootPath   string             `bson:"rootpath,omitempty"`
	ReleaseURL string             `bson:"releaseurl,omitempty"`
	Token      string             `bson:"token,omitempty"`
}

func Init(uri string) Config {
	i := Config{}
	if uri != "" {
		i.MongodbURI = uri
	}
	return i
}

func NewCollection(c db.Client) db.Collection {
	return c.NewDB("osoba").NewCollection("pathes", Document{})
}

func (c Config) Fetch(filter interface{}) ([]Document, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, c.MongodbURI)
	if err != nil {
		return []Document{}, err
	}
	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	pathes := client.NewDB("osoba").NewCollection("pathes", Document{})
	var docs []Document
	if err := pathes.Read(ctx, filter, &docs); err != nil {
		return []Document{}, err
	}

	pp.Println(docs)

	return docs, nil
}
