package repository

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// CreateUser handles creating new user.
func (mc *MongoDBClient) CreateUser(u *User) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("users")
	_, err := collection.InsertOne(ctx, u)
	// It may failed due to duplicate email as email has to be unique
	if err != nil {
		return fmt.Errorf("email already exists")
	}
	return nil
}

// UpdateProfileInformation updates the profile information
func (mc *MongoDBClient) UpdateProfileInformation(u *User) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("users")

	filterCondn := bson.M{
		"email": bson.M{"$regex": u.Email, "$options": "i"},
	}
	updateDoc := bson.M{"$set": bson.M{

		"firstName":     u.Firstname,
		"lastName":      u.LastName,
		"about":         u.About,
		"notifications": u.Notification,
		"resume":        u.Resume,
		"country":       u.Country,
	}}
	m, err := collection.UpdateOne(ctx, filterCondn, updateDoc)
	if m.MatchedCount == 0 {
		return fmt.Errorf("user does not exist.")
	}
	return err
}

func (mc *MongoDBClient) GetProfileByEmail(email string) (*User, error) {
	ctx, cancel := getContextWithTimeout(5)
	defer cancel()

	collection := mc.Database("referrer").Collection("users")
	m := collection.FindOne(ctx, bson.M{
		"email": bson.M{"$regex": email, "$options": "i"},
	})

	if m.Err() != nil {
		return nil, fmt.Errorf("no user found")
	}

	var u User
	err := m.Decode(&u)
	if err != nil {
		return nil, err
	}
	u.Password = ""

	return &u, nil
}
