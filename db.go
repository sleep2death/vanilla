package vanilla

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// DBName of the mongo
	DBName string = "vanilla"
	// UserCollection name
	UserCollection string = "users"
)

func initDB(addr string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		return nil, err
	}

	// test piing
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Println("connected to mongodb")
	return client.Database(DBName), nil
}
