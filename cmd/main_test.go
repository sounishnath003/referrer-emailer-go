package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestConnectMongoDB(t *testing.T) {

	mongoDbUri := utils.GetStringFromEnv("MONGO_DB_URI", "localhost")
	fmt.Printf("MongoDB Uri %s", mongoDbUri)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Ping the database to verify the connection is alive...
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		t.Errorf("Failed to ping MongoDB: %v", err)
	}
}
