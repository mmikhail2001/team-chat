package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Mongo *mongo.Database
}

func NewDatabase(MONGO_URI string, MONGO_DATABASE string) *Database {
	fmt.Println(MONGO_URI)
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	} else {
		log.Println("Successfully connected database")
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping database: %s", err)
	}

	// host > db > collection ???
	db := client.Database(MONGO_DATABASE)
	return &Database{
		Mongo: db,
	}
}
