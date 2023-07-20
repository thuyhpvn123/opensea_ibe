package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func InitDatabase() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://minigame:chinhtoi@cluster0.qpzfe.mongodb.net/?retryWrites=true&w=majority")
	client, _ := mongo.Connect(ctx, clientOptions)
	collection = client.Database("chiabai").Collection("chiabai")
	fmt.Println("Connection opened to database")

	return collection
}

// GetDB returns the DB instance.
func GetDB() *mongo.Collection {
	return collection
}
