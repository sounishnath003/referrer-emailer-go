package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Payload struct {
	ResumeURL        string `json:"resumeUrl" bson:"resumeUrl"`
	ExtractedContent string `json:"extractedContent" bson:"extractedContent"`
	Summary          string `json:"summary" bson:"summary"`
	ResumeJSON       string `json:"resumeJson" bson:"resumeJson"`
}

type JobQueue struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	UserEmailAddress string    `json:"userEmailAddress" bson:"userEmailAddress"`
	JobType          string    `json:"jobType" bson:"jobType"`
	Status           string    `json:"status" bson:"status"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
	Payload          Payload   `json:"payload" bson:"payload"`
}
