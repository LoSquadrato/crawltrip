package database

import (
	"log"

	"context"

	"github.com/LoSquadrato/crawltrip/internal/classificator"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TODO: after the broker implementation the database became a separate entity from the proxy.
// We don't need to store the database configuration inside the proxy

var databaseName = "crawltrip"
var collectionName = "requests"

func SaveRequest(ctx context.Context, client *mongo.Client, rawreq *classificator.RawRequest) error {
	log.Println("Saving request to database")
	// TODO: add option to save the request in the database, for now it is just a stub.
	res, err := client.Database(databaseName).Collection(collectionName).InsertOne(ctx, rawreq)
	if err != nil {
		return err
	}
	log.Printf("Request saved with id: %v\n", res.InsertedID)
	return nil
}

func DeleteAllRequests(ctx context.Context, client *mongo.Client) error {
	log.Println("Deleting all requests from database")
	_, err := client.Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	log.Println("All requests deleted from database")
	return nil
}

func DeleteRequestByID(ctx context.Context, client *mongo.Client, id string) error {
	log.Printf("Deleting request with id: %v from database\n", id)
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = client.Database(databaseName).Collection(collectionName).DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	log.Printf("Request with id: %v deleted from database\n", id)
	return nil
}

func GetRequestByID(ctx context.Context, client *mongo.Client, id string) (*classificator.RawRequest, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	log.Println("Getting request from database by ID")
	res := client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": objectId})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var rawreq classificator.RawRequest
	if err := res.Decode(&rawreq); err != nil {
		return nil, err
	}
	return &rawreq, nil
}

func GetLastRequest(ctx context.Context, client *mongo.Client) (*classificator.RawRequest, error) {
	log.Println("Getting last request from database")
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	res := client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{}, opts)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var rawreq classificator.RawRequest
	if err := res.Decode(&rawreq); err != nil {
		return nil, err
	}
	return &rawreq, nil
}
