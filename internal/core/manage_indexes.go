package core

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// getContextWithTimeout helps as utility to pass on the timeout controlled context.
func getContextWithTimeout(timeoutInSecond int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeoutInSecond)*time.Second)
}

// intiializeDatabase tries to connect with mongo DB database within 10 second context
func intiializeDatabase(dbUri string) (*mongo.Client, error) {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	slog.Info("mongo db client instance ping checked and connected...")
	return client, nil
}

func (co *Core) configureUsersIndexes() {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := co.DB.Database("referrer").Collection("users")

	foo, _ := bson.Marshal(bson.D{{Name: "email", Value: 1}})
	indexModel := mongo.IndexModel{
		Keys:    foo, // Create index on `email` field
		Options: options.Index().SetUnique(true).SetName("email_index"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		co.Lo.Error("failed to create index on users collection", slog.Any("index_err", err.Error()))
		panic(err)
	}

	co.Lo.Info("successfully created index on users collection")
}
