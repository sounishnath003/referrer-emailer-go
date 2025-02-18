package repository

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// GetProfileByEmail helps to find the user by email address
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
	if err != nil || len(u.Email) == 0 {
		return nil, fmt.Errorf("user or email is not found")
	}

	return &u, nil
}

// FindUserByEmailAndPassword user by email and password
func (mc *MongoDBClient) FindUserByEmailAndPassword(email, password string) (*User, error) {
	u, err := mc.GetProfileByEmail(email)

	if err != nil || u == nil || (u.Password != password) {
		return nil, fmt.Errorf("email or password is not valid")
	}

	// Create token.
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims.
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix() // 1 Hour

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte("Sec%&!*RT#*!@(89231%&!*RT#12345"))
	if err != nil {
		return nil, err
	}

	// Changing the password to empty, not to send in API response.
	u.Password = ""

	return u, nil

}
