package repository

import "gopkg.in/mgo.v2/bson"

type Skills struct {
	ProgrammingLanguages []string `json:"programmingLanguages" bson:"programmingLanguages"`
	ToolsAndTechnologies []string `json:"toolsAndTechnologies" bson:"toolsAndTechnologies"`
	Frameworks           []string `json:"frameworks" bson:"frameworks"`
	CloudPlatforms       []string `json:"cloudPlatforms" bson:"cloudPlatforms"`
	Miscellenous         []string `json:"miscellenous" bson:"miscellenous"`
}

type SocialLink struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

type WorkExperience struct {
	OrganizationName string `json:"organizationName" bson:"organizationName"`
	Location         string `json:"location" bson:"location"`
	Tenure           string `json:"tenure" bson:"tenure"`
	Experiences      string `json:"experiences" bson:"experiences"`
}

type PeronalProject struct {
	Name     string       `json:"name" bson:"name"`
	Links    []SocialLink `json:"projectDemos" bson:"projectDemos"`
	Features string       `json:"features" bson:"features"`
}

type Education struct {
	InstituteName string `json:"institutionName" bson:"institutionName"`
	MarksObtained string `json:"marksObtained" bson:"marksObtained"`
}
type Achievement struct {
	Detail string `json:"details" bson:"details"`
}

type ResumeInformation struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Email            string           `json:"email" bson:"email"`
	Skills           Skills           `json:"skills" bson:"skills"`
	SocialLinks      []SocialLink     `json:"socialLinks" bson:"socialLinks"`
	WorkExperiences  []WorkExperience `json:"workExperineces" bson:"workExperiences"`
	PersonalProjects []PeronalProject `json:"personalProjects" bson:"personalProjects"`
	Educations       []Education      `json:"educations" bson:"educations"`
	Achievements     []Achievement    `json:"achievements" bson:"achievements"`

	ProfileSummary string `json:"profileSummary" bson:"profileSummary"`
}
