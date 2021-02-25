package resource

import (
	"context"
	"log"
	"osoba/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (c Config) InitDB() {
	log.Println("DB initializing")
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, c.MongodbURI)
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

	return docs, nil
}
