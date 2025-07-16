package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TailoredResume struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"userId" json:"userId"`
	JobDescription string             `bson:"jobDescription" json:"jobDescription"`
	ResumeMarkdown string             `bson:"resumeMarkdown" json:"resumeMarkdown"`
	CompanyName    string             `bson:"companyName" json:"companyName"`
	JobRole        string             `bson:"jobRole" json:"jobRole"`
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

// UpdateTailoredResumeMarkdown updates the resumeMarkdown of a tailored resume by ID
func (mc *MongoDBClient) UpdateTailoredResumeMarkdown(ctx context.Context, id primitive.ObjectID, resumeMarkdown string) error {
	collection := mc.Database("referrer").Collection("tailored_resumes")
	_, err := collection.UpdateOne(ctx, primitive.M{"_id": id}, primitive.M{"$set": primitive.M{"resumeMarkdown": resumeMarkdown}})
	return err
}

// GetLatestTailoredResumesByUser fetches the latest 10 tailored resumes for a user by userId, optionally filtered by companyName
func (mc *MongoDBClient) GetLatestTailoredResumesByUser(ctx context.Context, userId string, companyName string) ([]*TailoredResume, error) {
	collection := mc.Database("referrer").Collection("tailored_resumes")
	filter := primitive.M{"userId": userId}
	if companyName != "" {
		filter["companyName"] = primitive.M{"$regex": companyName, "$options": "i"}
	}
	findOptions := options.Find().SetSort(primitive.M{"createdAt": -1}).SetLimit(10)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resumes []*TailoredResume
	for cursor.Next(ctx) {
		var tr TailoredResume
		if err := cursor.Decode(&tr); err != nil {
			return nil, err
		}
		resumes = append(resumes, &tr)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return resumes, nil
}
