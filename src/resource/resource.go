package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"osoba/db"
	"osoba/deploy"
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

	return docs, nil
}

func (c Config) FetchAll() ([]Document, error) {
	docs, err := c.Fetch(bson.M{})
	if err != nil {
		return []Document{}, err
	}
	return docs, nil
}

func (c Config) FetchJsonAll() ([]byte, error) {
	docs, err := c.Fetch(bson.M{})
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(docs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DocToJsonWithoutToken(docs []Document) ([]byte, error) {
	type document struct {
		Path       string
		ReleaseURL string
	}

	d := []document{}
	for _, doc := range docs {
		d = append(d, document{
			Path:       doc.Path,
			ReleaseURL: doc.ReleaseURL,
		})
	}

	res, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c Config) FetchJsonAllWithoutToken() ([]byte, error) {
	docs, err := c.FetchAll()
	if err != nil {
		return nil, err
	}

	res, err := DocToJsonWithoutToken(docs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func MapToDoc(m map[string]string) Document {
	doc := Document{
		Path:       m["path"],
		ReleaseURL: m["releaseurl"],
		Token:      m["token"],
	}
	log.Println(doc)
	return doc
}

func (c Config) SetDoc(doc Document) error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, c.MongodbURI)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	pathes := client.NewDB("osoba").NewCollection("pathes", Document{})
	var docs []Document
	if err := pathes.Read(ctx, bson.M{"path": doc.Path}, &docs); err != nil {
		return err
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(doc.Token), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	doc.Token = string(hashedToken)

	log.Printf("%#v\n", doc)

	if len(docs) == 0 {
		return pathes.Insert(ctx, []Document{doc})
	} else if len(docs) == 1 {
		return pathes.Update(ctx, bson.M{"path": doc.Path},
			bson.D{
				{"$set", bson.D{{"field1", "xxx"}}},
			})
	}

	return fmt.Errorf("docs length is invalid, path: %s, docs: %#v", doc.Path, docs)
}

func (c Config) Delete(path string) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, c.MongodbURI)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	pathes := client.NewDB("osoba").NewCollection("pathes", Document{})
	var docs []Document
	if err := pathes.Read(ctx, bson.M{"path": path}, &docs); err != nil {
		return nil, err
	}

	if len(docs) != 1 {
		return nil, fmt.Errorf("docs length is invalid, path: %s, docs: %#v", path, docs)
	}

	if err := pathes.Delete(ctx, bson.M{"path": path}); err != nil {
		return nil, err
	}

	res, err := DocToJsonWithoutToken(docs)
	if err != nil {
		return nil, err
	}

	// TODO to be async
	di := deploy.Info{
		Path: docs[0].Path,
	}
	if err := di.Delete(); err != nil {
		return nil, err
	}

	return res, nil
}
