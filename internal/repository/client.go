package repository

import "go.mongodb.org/mongo-driver/mongo"

type MongoDBClient struct {
	*mongo.Client
}
