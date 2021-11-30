package person

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultDB      = "persondb"
	collectionName = "person"
)

type MongoHandler struct {
	client   *mongo.Client
	database string
}

// constructor
func NewHandler(address string) *MongoHandler {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	mongoHandler := &MongoHandler{
		client:   client,
		database: defaultDB,
	}
	return mongoHandler
}

func (mongoHandler *MongoHandler) GetOne(person *Person, filter interface{}) error {
	// will create the collection if none
	collection := mongoHandler.client.Database(mongoHandler.database).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(person)
	return err
}

func (mongoHandler *MongoHandler) Get(filter interface{}) []*Person {
	collection := mongoHandler.client.Database(mongoHandler.database).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []*Person
	for cursor.Next(ctx) {
		person := &Person{}
		err1 := cursor.Decode(person)
		if err1 != nil {
			log.Fatal(err1)
		}
		result = append(result, person)
	}
	return result
}

func (mongoHandler *MongoHandler) AddOne(person *Person) (*mongo.InsertOneResult, error) {
	collection := mongoHandler.client.Database(mongoHandler.database).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, person)
	return result, err
}

func (mongoHandler *MongoHandler) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mongoHandler.client.Database(mongoHandler.database).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.UpdateMany(ctx, filter, update)
	return result, err
}

func (MongoHandler *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := MongoHandler.client.Database(MongoHandler.database).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
