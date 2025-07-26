package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AiDraftColdEmail struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserEmailAddress string        `json:"userEmailAddress,omitempty" bson:"userEmailAddress"`

	To   string `json:"to,omitempty" bson:"to"`
	From string `json:"from,omitempty" bson:"from"`

	CompanyName    string   `json:"companyName,omitempty" bson:"companyName"`
	JobUrls        []string `json:"jobUrls,omitempty" bson:"jobUrls"`
	JobDescription string   `json:"jobDescription,omitempty" bson:"jobDescription"`
	TemplateType   string   `json:"templateType,omitempty" bson:"templateType"`

	MailSubject string `json:"mailSubject,omitempty" bson:"mailSubject"`
	Mailbody    string `json:"mailBody,omitempty" bson:"mailBody"`

	TailoredResumeID string `json:"tailoredResumeId,omitempty" bson:"tailoredResumeId"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type ProfileAnalytics struct {
	TotalEmails int                     `json:"totalEmails" bson:"totalEmails"`
	Companies   []CompanyEmailAggregate `json:"companies" bson:"companies"`
}

// CompanyEmailAggregate represents company-wise aggregation
type CompanyEmailAggregate struct {
	CompanyName        string `json:"companyName" bson:"companyName"`
	TotalEmails        int    `json:"totalEmails" bson:"totalEmails"`
	DistinctUsersCount int    `json:"distinctUsersCount" bson:"distinctUsersCount"`
}
