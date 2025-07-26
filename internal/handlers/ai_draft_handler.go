package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

type ReferralColdmailRequestDto struct {
	To   string `json:"to"`
	From string `json:"from"`

	CompanyName    string   `json:"companyName"`
	JobUrls        []string `json:"jobUrls"`
	JobDescription string   `json:"jobDescription"`
	TemplateType   string   `json:"templateType"`
}

func DraftReferralEmailWithAiHandler(c echo.Context) error {
	// Get context
	hctx := c.(*HandlerContext)

	var rmailDto ReferralColdmailRequestDto

	if err := c.Bind(&rmailDto); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// check for errors of proper email address
	if !isValidEmail(rmailDto.To) || !isValidEmail(rmailDto.From) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email address"))
	}

	// Steps to perform the action
	// Step01: Get the profile information from the `from` email address.
	u, err := hctx.GetCore().DB.GetProfileByEmail(rmailDto.From)
	if err != nil || len(u.Firstname) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unable to fetch information for %s: %w", rmailDto.From, err))
	}
	// Step02: Call LLM apis with private customized prompt.
	mailSubject, mailBody, err := hctx.GetCore().DraftColdEmailMessageLLM(
		rmailDto.From,
		rmailDto.To,
		rmailDto.CompanyName,
		rmailDto.TemplateType,
		rmailDto.JobDescription,
		u.ProfileSummary,
		rmailDto.JobUrls,
	)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unable to generated draft email %s: %w", rmailDto.From, err))
	}

	var tailoredResumeID string

	// If Jobdescription has been provided, then only generate a tailored resume
	if len(rmailDto.JobDescription) > 0 {
		// Step03: Generate tailored resume and store it
		resumeMarkdown, tErr := hctx.GetCore().TailorResumeWithJobDescriptionLLM(rmailDto.JobDescription, u.ExtractedContent, rmailDto.CompanyName, rmailDto.TemplateType)
		if tErr == nil {
			tr := &repository.TailoredResume{
				UserID:         u.ID.Hex(),
				JobDescription: rmailDto.JobDescription,
				ResumeMarkdown: resumeMarkdown,
				CompanyName:    rmailDto.CompanyName,
				JobRole:        rmailDto.TemplateType,
			}
			ctx := c.Request().Context()
			insertedID, err := hctx.GetCore().DB.CreateTailoredResume(ctx, tr)
			if err == nil {
				tailoredResumeID = insertedID.Hex()
			}
		}
	}

	// Step04: Store the draft email with tailoredResumeID
	_, _ = hctx.GetCore().DB.CreateAiDraftEmail(
		rmailDto.From,
		rmailDto.To,
		rmailDto.CompanyName,
		rmailDto.TemplateType,
		rmailDto.JobDescription,
		u.ProfileSummary,
		mailSubject,
		mailBody,
		tailoredResumeID,
		rmailDto.JobUrls,
	)

	return c.JSON(http.StatusOK, map[string]string{
		"mailSubject":      mailSubject,
		"mailBody":         mailBody,
		"tailoredResumeId": tailoredResumeID,
	})
}
