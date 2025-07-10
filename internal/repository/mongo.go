package repository

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
)

const mongoCaller = packageCaller + ".Mongo"
const _ = mongoCaller

type Mongo struct {
	logger *utils.Logger
	client *mongo.Client
}

func NewMongoRepo(logger *utils.Logger) ports.MongoRepo {
	host := os.Getenv("MONGO_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "MONGO_HOST is not set")
	}
	port := os.Getenv("MONGO_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "MONGO_PORT is not set")
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://" + host + ":" + port + "/?replicaSet=rs0"),
	)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to ping MongoDB: %v", err)
	}
	return &Mongo{
		logger: logger,
		client: client,
	}
}

func (r *Mongo) Close() error {
	if err := r.client.Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}
