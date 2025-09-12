package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	connected  bool
}

func NewMongoDBRepository(connectionString string) (*MongoDBRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Testar a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	database := client.Database("userdb")
	collection := database.Collection("processed_users")

	log.Println("Successfully connected to MongoDB")

	return &MongoDBRepository{
		client:     client,
		database:   database,
		collection: collection,
		connected:  true,
	}, nil
}

func (r *MongoDBRepository) SaveProcessedData(data *ProcessedData) (string, error) {
	if !r.connected {
		return "", fmt.Errorf("MongoDB not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

func (r *MongoDBRepository) GetAllProcessedData() ([]ProcessedData, error) {
	if !r.connected {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ProcessedData
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *MongoDBRepository) IsConnected() bool {
	return r.connected
}

func (r *MongoDBRepository) Close() error {
	if r.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return r.client.Disconnect(ctx)
	}
	return nil
}
