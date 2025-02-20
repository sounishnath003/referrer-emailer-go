package repository

import "gopkg.in/mgo.v2/bson"

// Notification reference model defination for the notification model
type Notification struct {
	Offers           bool   `json:"offers" bson:"offers"`
	PushNotification string `json:"pushNotifications" bson:"pushNotifications"`
	ReceiveEmails    bool   `json:"receiveEmails" bson:"receiveEmails"`
}

// User model to construct the mongoDB user model
type User struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Firstname    string       `json:"firstName" bson:"firstName"`
	LastName     string       `json:"lastName" bson:"lastName"`
	Resume       string       `json:"resume" bson:"resume"`
	About        string       `json:"about" bson:"about"`
	Country      string       `json:"country" bson:"country"`
	Notification Notification `json:"notifications" bson:"notifications"`

	Email    string `json:"email" bson:"email"`
	Password string `json:"password,omitempty" bson:"password"`

	ProfileSummary   string `json:"profileSummary" bson:"profileSummary"`
	ExtractedContent string `json:"extractedContent" bson:"extractedContent"`

	Token string `json:"token,omitempty" bson:"-"`
}
