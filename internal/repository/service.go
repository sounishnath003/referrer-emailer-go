package repository

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		"firstName":        u.Firstname,
		"lastName":         u.LastName,
		"about":            u.About,
		"notifications":    u.Notification,
		"resume":           u.Resume,
		"country":          u.Country,
		"profileSummary":   u.ProfileSummary,
		"extractedContent": u.ExtractedContent,
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
	claims["subject"] = u.Email
	claims["email"] = u.Email
	claims["iss"] = "referrer-emailer-service"
	claims["iat"] = time.Now().Unix()
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

// CreateEmailInMailbox stores the email into mailbox
func (mc *MongoDBClient) CreateEmailInMailbox(from string, to []string, subject, body string) error {
	// Get context
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	mail := ReferralMailbox{
		Uuid:      uuid.New().String(),
		From:      from,
		To:        to,
		Subject:   subject,
		Body:      body,
		CreatedAt: time.Now(),
	}

	collection := mc.Database("referrer").Collection("referral_mailbox")
	_, err := collection.InsertOne(ctx, mail)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoDBClient) GetLatestEmailsByFilter(filterCondn bson.M) ([]*ReferralMailbox, error) {
	// Get context.
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	var mails []*ReferralMailbox

	collection := mc.Database("referrer").Collection("referral_mailbox")
	fmt.Println(filterCondn)

	cursor, err := collection.Find(ctx, filterCondn, options.Find().SetLimit(10).SetSort(bson.M{"createdAt": -1}))
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	if err := cursor.Err(); err != nil {
		return mails, nil
	}

	for cursor.Next(ctx) {
		var r ReferralMailbox
		err := cursor.Decode(&r)
		if err != nil {
			fmt.Println("error occured", err)
			return nil, err
		}
		mails = append(mails, &r)
	}

	if len(mails) == 0 {
		return mails, fmt.Errorf("no results found")
	}

	return mails, nil
}

func (mc *MongoDBClient) CreateAiDraftEmail(from, to, companyName, templateType, jobDescription, userProfileSummary, mailSubject, mailBody string, jobUrls []string) (*AiDraftColdEmail, error) {
	// Get context
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	draftEmail := &AiDraftColdEmail{
		UserEmailAddress: from,
		From:             from,
		To:               to,
		CompanyName:      companyName,
		JobUrls:          jobUrls,
		JobDescription:   jobDescription,
		TemplateType:     templateType,
		MailSubject:      mailSubject,
		Mailbody:         mailBody,
		CreatedAt:        time.Now(),
	}

	collection := mc.Database("referrer").Collection("ai_email_drafts")
	_, err := collection.InsertOne(ctx, draftEmail)
	if err != nil {
		return nil, err
	}

	return draftEmail, nil
}
