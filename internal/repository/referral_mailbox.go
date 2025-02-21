package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ReferralMailbox struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Uuid      string        `json:"uuid" bson:"uuid"`
	From      string        `json:"from" bson:"from"`
	To        []string      `json:"to" bson:"to"`
	Subject   string        `json:"subject" bson:"subject"`
	Body      string        `json:"body" bson:"body"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
}
