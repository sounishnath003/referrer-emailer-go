package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BulkEmailJob struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserEmail       string             `json:"userEmail" bson:"userEmail"`
	TotalRecipients int                `json:"totalRecipients" bson:"totalRecipients"`
	SentCount       int                `json:"sentCount" bson:"sentCount"`
	Status          string             `json:"status" bson:"status"` // PENDING, IN_PROGRESS, COMPLETED, FAILED
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	Errors          []string           `json:"errors" bson:"errors"`
}

func (mc *MongoDBClient) CreateBulkEmailJob(job *BulkEmailJob) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()
	
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	
	collection := mc.Database("referrer").Collection("bulk_email_jobs")
	result, err := collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}
	job.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (mc *MongoDBClient) UpdateBulkEmailJobProgress(id primitive.ObjectID, sentCount int) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("bulk_email_jobs")
	
	update := bson.M{
		"$set": bson.M{
			"sentCount": sentCount,
			"updatedAt": time.Now(),
			"status":    "IN_PROGRESS",
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (mc *MongoDBClient) CompleteBulkEmailJob(id primitive.ObjectID) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("bulk_email_jobs")
	
	update := bson.M{
		"$set": bson.M{
			"status":    "COMPLETED",
			"updatedAt": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (mc *MongoDBClient) FailBulkEmailJob(id primitive.ObjectID, errStr string) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("bulk_email_jobs")
	
	update := bson.M{
		"$set": bson.M{
			"status":    "FAILED",
			"updatedAt": time.Now(),
		},
		"$push": bson.M{
			"errors": errStr,
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (mc *MongoDBClient) GetBulkEmailJob(id primitive.ObjectID) (*BulkEmailJob, error) {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("bulk_email_jobs")
	
	var job BulkEmailJob
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
