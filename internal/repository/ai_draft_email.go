package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AiDraftColdEmail struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserEmailAddress string        `json:"userEmailAddress,omitempty" bson:"userEmailAddress" bson:"userEmailAddress"`

	To   string `json:"to,omitempty" bson:"to"`
	From string `json:"from,omitempty" bson:"from"`

	CompanyName    string   `json:"companyName,omitempty" bson:"companyName"`
	JobUrls        []string `json:"jobUrls,omitempty" bson:"jobUrls"`
	JobDescription string   `json:"jobDescription,omitempty" bson:"jobDescription"`
	TemplateType   string   `json:"templateType,omitempty" bson:"templateType"`

	MailSubject string `json:"mailSubject,omitempty" bson:"mailSubject"`
	Mailbody    string `json:"mailBody,omitempty" bson:"mailBody"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
