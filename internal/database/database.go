// This package provides functions to connect to a MongoDB database,
// ping the database to check the connection,
// and close the connection when done.
// It uses the official MongoDB Go driver to manage the connection and context.
package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Close(client *mongo.Client, ctx context.Context) {
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			fmt.Println("Error disconnecting from database:", err)
		}
	}()
}

func Connect(uri string, timeout time.Duration) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	opts.SetServerSelectionTimeout(timeout)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
