package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"api-rabbitmq/internal/domain/entities"
	"api-rabbitmq/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryImpl struct {
	client     *mongo.Client
	collection *mongo.Collection
	connected  bool
}

func NewUserRepository(connectionString string) (repositories.UserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	database := client.Database("userdb")
	collection := database.Collection("processed_users")

	log.Println("Successfully connected to MongoDB")

	return &UserRepositoryImpl{
		client:     client,
		collection: collection,
		connected:  true,
	}, nil
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *entities.ProcessedUser) (string, error) {
	user.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]entities.ProcessedUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []entities.ProcessedUser
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.ProcessedUser, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID: %v", err)
	}

	var user entities.ProcessedUser
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Close() error {
	if r.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return r.client.Disconnect(ctx)
	}
	return nil
}
