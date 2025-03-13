package core

import (
	"errors"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

// initializeLLM initializes the LLM (Large Language Model) client for the Core instance.
// It sets up a context with a timeout of 5 seconds and creates a new Generative AI client
// using the provided GCP project ID and location. If the client creation is successful,
// it assigns the generative model to the Core instance. The client is closed before the
// function returns.
//
// Returns an error if the client creation fails.
func (co *Core) initializeLLM() error {
	ctx, cancel := getContextWithTimeout(5)
	defer cancel()

	client, err := genai.NewClient(ctx, co.opts.GcpProjectID, co.opts.GcpLocation)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	co.llm = client.GenerativeModel(co.opts.ModelName)

	return nil
}

func (co *Core) ExtractResumeContentLLM(resumePath string) (string, error) {
	co.Lo.Info("started extracting content", "resume", resumePath)

	ctx, cancel := getContextWithTimeout(30)
	defer cancel()

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  resumePath,
	}

	res, err := co.llm.GenerateContent(ctx, part, genai.Text(`
		[Backstory]: You are a very "Professional Hiring Manager Specialist" of a Company. You have to screen candidate based on their resume.

		[Task]: 
		- Given a "Resume", You have to extract all the "Key Informations" of this candidate. 
		- Start with "## Candidate Name" . DO NOT ADD ANYTHING ELSE.

		[Important to note]:
    	1. Extract all the Social Links (Linkedin, Github, MobileNo, EmailAdress etc).
		2. Extract the "Skillsets" which are important for a Job description to match.
		3. Extract all "Work Experiences" of this candidate in Chronological Order. Along with the "Key Performance" highlights.
		4. Extract all "Projects / Personal Projects" candidate has mentioned. Capture the "Technologies has been used", "Project Demo links" (Github, Youtube, Or similar) etc.
		5. Extract "Achievements" of this candidates if present. Else "SKIP".
		6. Extract "Latest Education" (College, Institutions), Year of Graduation Tenure. "SKIP" the Class "10th, 12th" portions.
	`))
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	llmContent := fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])

	return llmContent, nil
}

// GenerateProfileSummaryLLM generates a summary of a resume using an AI model.
// It takes the path to the resume file as input and logs the process.
// The function uploads the resume to a cloud storage and sends it to the AI model for summarization.
// If the AI model returns a valid summary, it logs the summary.
// If there is an error during the process or the AI model returns an empty response, it returns an error.
//
// Parameters:
//   - resumePath: The path to the resume file to be summarized.
//
// Returns:
//   - error: An error if the summarization process fails or the AI model returns an empty response.
func (co *Core) GenerateProfileSummaryLLM(content string) (string, error) {
	ctx, cancel := getContextWithTimeout(30)
	defer cancel()

	res, err := co.llm.GenerateContent(ctx, genai.Text(content), genai.Text(`
			[Backstory]: You are a very "Professional Content Generation Manager Specialist".
			[Task]: Given a "Resume Content", Please summarize the content into "BULLET POINTS".
			
			[Important to note]:
			1. Give more attention to the "Professional Work, Skills, Project and Achievements".
			2. Must Keep the Contact Details (phone, email, linkedin, portfolio, etc) in your summary.
			3. Output must be "Markdown" format
	`))
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	llmContent := fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
	return llmContent, nil
}

func (co *Core) ConvertResumeToJSONStructLLM(content string) (string, error) {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	res, err := co.llm.GenerateContent(ctx, genai.Text(content), genai.Text(`
		[Backstory]: You are a "Senior Data Entry Specialist". Your task is to MAP the "Content" into "JSON" structure.
		[Task]: Given a "Resume" content, Generate output Provided Below - 

		'{"fullName":"string","email":"string","skills":{"programmingLanguages":["string"],"toolsAndTechnologies":["string"],"frameworks":["string"],"cloudPlatforms":["string"],"miscellenous":["string"]},"socialLinks":[{"type":"string","value":"string"}],"workExperiences":[{"organizationName":"string","location":"string","tenure":"string","experiences":"string"}],"personalProjects":[{"name":"string","projectDemos":[{"type":"string","value":"string"}],"features":"string"}],"educations":[{"institutionName":"string","marksObtained":"string"}],"achievements":[{"details":"string"}]}'

	`))

	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	llmContent := fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
	return llmContent, nil
}

// DraftColdEmailMessageLLM helps to generate a draft email and returns (mailSubject, mailBody, error).
//
// It also stores all the Ai generated draft in database to track.
func (co *Core) DraftColdEmailMessageLLM(from, to, companyName, templateType, jobDescription, userProfileSummary string, jobUrls []string) (string, string, error) {
	ctx, cancel := getContextWithTimeout(15)
	defer cancel()

	res, err := co.llm.GenerateContent(
		ctx,
		genai.Text(
			fmt.Sprintf(`
				** JOB Opportunity Details:**

				To: %s
				CompanyName: %s
				JOB URLs: %v,
				JobDescription: %s
				
				** Candidate Profile:**

				%s

			`, to, companyName, jobUrls, jobDescription, userProfileSummary),
		),
		genai.Text(fmt.Sprintf(`		
			I am a [Your Role] with [Years of Experience] experience in [Industry/Area]. I'm writing to express my interest in [Job Title] at [Company Name]. I found the job posting on [Platform]. Write a cold email to the recruiter [ToEmail]. Highlight my relevant skills and experience and REQUESTING to SCHEDULE AN INTERVIEW!. with STAR method like "Performed X with Y and achieved Z%".
			
			Keep it under 200 words. Write it in "1st Person Candidate's View". While adding "JOB URLs add in Bullet list" manner.

			**Specific Requirements:**

			1. Use more Bullet Points and Bold Keywords.
			2.  Include a candidate signature (Contact Details: (phone, email, linkedin, portfolio, etc), utilizing information from the "Candidate Profile."
			3.  Format the entire output as "Markdown" format. 
		`, templateType)),
	)

	if err != nil {
		return "", "", fmt.Errorf("unable to generate mailbody contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", "", errors.New("empty response mailBody from model")
	}

	mailBody := fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])

	res, err = co.llm.GenerateContent(ctx, genai.Text(mailBody), genai.Text(`
		[Task]: Given "Referral Cold Email Body" in "HTML" format, targets for Which Type of Job (Software Engineering, Backend Developer, Data Engineering, Data Analyst, etc).
		Reply In One Word "TYPE OF THE JOB".
	`))
	if err != nil {
		return "", "", fmt.Errorf("unable to generate type.Of.Job contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", "", errors.New("empty response type.Of.Job from model")
	}

	mailSubject := fmt.Sprintf("Interested for %s Roles Opportunities at %s", res.Candidates[0].Content.Parts[0], companyName)

	// mailSubject := strings.Split(mailBody, "\n\n")[0]

	// Store the draft into DB
	_, err = co.DB.CreateAiDraftEmail(from, to, companyName, templateType, jobDescription, userProfileSummary, mailSubject, mailSubject, jobUrls)

	return mailSubject, mailBody, nil
}
