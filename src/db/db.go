package db

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
}

func NewClient(ctx context.Context, URI string) (Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	return Client{client}, err
}

func (c Client) DisconnectDB(ctx context.Context) error {
	return c.Disconnect(ctx)
}

func (c Client) NewDB(name string) DB {
	return DB{c.Database(name)}
}

type DB struct {
	*mongo.Database
}

func (db DB) NewCollection(name string, doc interface{}) Collection {
	return Collection{db.Collection(name), reflect.TypeOf(doc)}
}

type Collection struct {
	*mongo.Collection
	docType reflect.Type
}

func (c Collection) Insert(ctx context.Context, docs interface{}) error {
	if reflect.TypeOf(docs) != reflect.SliceOf(c.docType) {
		return fmt.Errorf("Error: type of docs is invalid, %#v\n", docs)
	}
	documents := []interface{}{}
	v := reflect.ValueOf(docs).Convert(reflect.SliceOf(c.docType))
	for i := 0; i < v.Len(); i++ {
		documents = append(documents, v.Index(i).Interface())
	}
	_, err := c.InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}

func (c Collection) Read(ctx context.Context, filter interface{}, docs interface{}) error {
	if reflect.TypeOf(docs) != reflect.PtrTo(reflect.SliceOf(c.docType)) {
		return fmt.Errorf("Error: type of docs is invalid, %#v\n", docs)
	}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, docs); err != nil {
		return err
	}
	return nil
}
func (c Collection) Update(ctx context.Context, filter, update interface{}) error {
	_, err := c.UpdateMany(ctx, filter, update)
	return err
}

func (c Collection) Delete(ctx context.Context, filter interface{}) error {
	_, err := c.DeleteMany(ctx, filter)
	return err
}
