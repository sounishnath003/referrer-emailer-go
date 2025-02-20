package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type JobType int

const (
	EXTRACT_CONTENT = iota + 1
	GENERATE_PROFILE_SUMMARY
	UPDATE_RESUME_DOCUMENT
	EMAIL_NOTIFICATION
)

// String - Creating common behavior - give the type a String function
func (j JobType) String() string {
	return [...]string{"EXTRACT_CONTENT", "GENERATE_PROFILE_SUMMARY", "UPDATE_RESUME_DOCUMENT", "EMAIL_NOTIFICATION"}[j-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (j JobType) EnumIndex() int {
	return int(j)
}

type Payload struct {
	ResumeURL        string `json:"resumeUrl" bson:"resumeUrl"`
	ExtractedContent string `json:"extractedContent" bson:"extractedContent"`
	Summary          string `json:"summary" bson:"summary"`
}

type JobQueue struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	UserEmailAddress string    `json:"userEmailAddress" bson:"userEmailAddress"`
	JobType          JobType   `json:"jobType" bson:"jobType"`
	Status           string    `json:"status" bson:"status"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
	Payload          Payload   `json:"payload" bson:"payload"`
}
