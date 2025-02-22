package core

// DraftColdEmailAi represents the structure of a cold email drafted by AI.
type DraftColdEmailAi struct {
	From           string
	To             string
	JobUrls        []string
	JobDescription string
	TemplateType   string
	Subject        string
	Body           string
}

// DraftColdEmailBuilder defines the interface for building a DraftColdEmailAi.
type DraftColdEmailBuilder interface {
	SetFrom(string) DraftColdEmailBuilder
	SetTo(string) DraftColdEmailBuilder
	SetJobUrls([]string) DraftColdEmailBuilder
	SetJobDescription(string) DraftColdEmailBuilder
	SetTemplateType(string) DraftColdEmailBuilder
	SetSubject(string) DraftColdEmailBuilder
	SetBody(string) DraftColdEmailBuilder

	Build() *DraftColdEmailAi
}

// NewDraftColdEmailBuilder creates a new instance of draftColdEmailBuilder.
func NewDraftColdEmailBuilder() DraftColdEmailBuilder {
	return &draftColdEmailBuilder{
		draftColdEmail: &DraftColdEmailAi{},
	}
}

// draftColdEmailBuilder is the concrete implementation of DraftColdEmailBuilder.
type draftColdEmailBuilder struct {
	draftColdEmail *DraftColdEmailAi
}

func (dc *draftColdEmailBuilder) SetFrom(fromEmail string) DraftColdEmailBuilder {
	dc.draftColdEmail.From = fromEmail
	return dc
}

func (dc *draftColdEmailBuilder) SetTo(toEmail string) DraftColdEmailBuilder {
	dc.draftColdEmail.To = toEmail
	return dc
}

func (dc *draftColdEmailBuilder) SetJobUrls(jobUrls []string) DraftColdEmailBuilder {
	dc.draftColdEmail.JobUrls = jobUrls
	return dc
}

func (dc *draftColdEmailBuilder) SetJobDescription(jobDesc string) DraftColdEmailBuilder {
	dc.draftColdEmail.JobDescription = jobDesc
	return dc
}

func (dc *draftColdEmailBuilder) SetTemplateType(templateType string) DraftColdEmailBuilder {
	dc.draftColdEmail.TemplateType = templateType
	return dc
}

func (dc *draftColdEmailBuilder) SetSubject(subject string) DraftColdEmailBuilder {
	dc.draftColdEmail.Subject = subject
	return dc
}

func (dc *draftColdEmailBuilder) SetBody(body string) DraftColdEmailBuilder {
	dc.draftColdEmail.Body = body
	return dc
}

func (dc *draftColdEmailBuilder) Build() *DraftColdEmailAi {
	return dc.draftColdEmail
}
