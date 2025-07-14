package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TailoredResume struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"userId" json:"userId"`
	JobDescription string             `bson:"jobDescription" json:"jobDescription"`
	ResumeMarkdown string             `bson:"resumeMarkdown" json:"resumeMarkdown"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	// Add more metadata fields as needed
}

// CreateTailoredResume stores a tailored resume in the tailored-resume collection
func (mc *MongoDBClient) CreateTailoredResume(ctx context.Context, tr *TailoredResume) (primitive.ObjectID, error) {
	tr.CreatedAt = time.Now()
	collection := mc.Database("referrer").Collection("tailored_resumes")
	result, err := collection.InsertOne(ctx, tr)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetTailoredResumeByID fetches a tailored resume by its ID
func (mc *MongoDBClient) GetTailoredResumeByID(ctx context.Context, id primitive.ObjectID) (*TailoredResume, error) {
	collection := mc.Database("referrer").Collection("tailored_resumes")
	var tr TailoredResume
	err := collection.FindOne(ctx, primitive.M{"_id": id}).Decode(&tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}
