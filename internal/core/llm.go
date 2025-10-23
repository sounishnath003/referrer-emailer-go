package core

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/genai"
)

// initializeLLM initializes the LLM (Large Language Model) client for the Core instance.
func (co *Core) initializeLLM() error {
	ctx, cancel := getContextWithTimeout(5)
	defer cancel()

	// Correct client initialization for Vertex AI
	timeDuration := 60 * time.Second
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			APIVersion: "v1",
			Timeout:    &timeDuration,
			Headers: http.Header{
				"X-Vertex-AI-LLM-Request-Type": []string{"shared"},
			},
		},
		Project: co.opts.GcpProjectID, Location: co.opts.GcpLocation, Backend: genai.BackendVertexAI,},
	)
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}
	co.llm = client

	return nil
}

func printResponse(res *genai.GenerateContentResponse) string {
	var result string
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				result += fmt.Sprint(part)
			}
		}
	}
	return result
}

func (co *Core) ExtractResumeContentLLM(resumePath string) (string, error) {
	co.Lo.Info("started extracting content", "resume", resumePath)

	ctx, cancel := getContextWithTimeout(30)
	defer cancel()

	res, err := co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{FileData: &genai.FileData{MIMEType: "application/pdf", FileURI: resumePath}},
			{Text: `
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
	`},
		},
	}}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", errors.New("empty response from model")
	}

	return printResponse(res), nil
}

// GenerateProfileSummaryLLM generates a summary of a resume using an AI model.
func (co *Core) GenerateProfileSummaryLLM(content string) (string, error) {
	ctx, cancel := getContextWithTimeout(30)
	defer cancel()

	res, err := co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{Text: content},
			{Text: `
			[Backstory]: You are a very "Professional Content Generation Manager Specialist".
			[Task]: Given a "Resume Content", Please summarize the content into "BULLET POINTS".
			
			[Important to note]:
			1. Give more attention to the "Professional Work, Skills, Project and Achievements".
			2. Must Keep the Contact Details (phone, email, linkedin, portfolio, etc) in your summary.
			3. Output must be "Markdown" format
	`},
		},
	}}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", errors.New("empty response from model")
	}

	return printResponse(res), nil
}

func (co *Core) ConvertResumeToJSONStructLLM(content string) (string, error) {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	res, err := co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{Text: content},
			{Text: `
		[Backstory]: You are a "Senior Data Entry Specialist". Your task is to MAP the "Content" into "JSON" structure.
		[Task]: Given a "Resume" content, Generate output Provided Below - 

		'{"fullName":"string","email":"string","skills":{"programmingLanguages":["string"],"toolsAndTechnologies":["string"],"frameworks":["string"],"cloudPlatforms":["string"],"miscellenous":["string"]},"socialLinks":[{"type":"string","value":"string"}],"workExperiences":[{"organizationName":"string","location":"string","tenure":"string","experiences":"string"}],"personalProjects":[{"name":"string","projectDemos":[{"type":"string","value":"string"}],"features":"string"}],"educations":[{"institutionName":"string","marksObtained":"string"}],"achievements":[{"details":"string"}]}'

	`},
		},
	}}, nil)

	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", errors.New("empty response from model")
	}

	return printResponse(res), nil
}

// DraftColdEmailMessageLLM helps to generate a draft email and returns (mailSubject, mailBody, error).
func (co *Core) DraftColdEmailMessageLLM(from, to, companyName, templateType, jobDescription, userProfileSummary string, jobUrls []string) (string, string, error) {
	ctx, cancel := getContextWithTimeout(15)
	defer cancel()

	res, err := co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{Text: fmt.Sprintf(`
				** JOB Opportunity Details:**

				To: %s
				CompanyName: %s
				JOB URLs: %v,
				JobDescription: %s
				
				** Candidate Profile:**

				%s

			`, to, companyName, jobUrls, jobDescription, userProfileSummary)},
			{Text: fmt.Sprintf(`		
			Write a cold email to the recruiter [ToEmail]. Highlight my relevant skills and experience and REQUESTING to SCHEDULE AN INTERVIEW!. with STAR method like "Performed X with Y and achieved Z%%".
			
			Keep it under 200 words. Write it in "1st Person Candidate's View". While adding "JOB URLs add in Bullet list" manner.

			**Specific Requirements:**

			1. Use more Bullet Points and Bold Keywords.
			2.  Include a candidate signature (Contact Details: (phone, email, linkedin, portfolio, etc), utilizing information from the "Candidate Profile."
			3.  Format the entire output as "Markdown" format. 
			4. No Need of providing a Subject Line.
		`, templateType)},
		},
	}}, nil)

	if err != nil {
		return "", "", fmt.Errorf("unable to generate mailbody contents: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", "", errors.New("empty response mailBody from model")
	}

	mailBody := printResponse(res)

	res, err = co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{Text: mailBody},
			{Text: fmt.Sprintf(`
		[Task]: Based on the "Referral Cold Email Body" and the job description, and the company name "%s", generate a concise, professional, and attention-grabbing email subject line (Mail Heading) that would make a recruiter want to open the email. The subject should clearly reflect the candidate's intent and the job opportunity, and be tailored to the content of the cold email and the job/company context.  The subject line MUST explicitly mention the company name ("%s").
		Reply ONLY with the subject line, no extra commentary or formatting.
	`, companyName, companyName)},
		},
	}}, nil)
	if err != nil {
		return "", "", fmt.Errorf("unable to generate type.Of.Job contents: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", "", errors.New("empty response type.Of.Job from model")
	}

	mailSubjectContent := printResponse(res)
	mailSubject := fmt.Sprintf("Interested for %s - %s", mailSubjectContent, companyName)

	// Store the draft into DB
	_, err = co.DB.CreateAiDraftEmail(from, to, companyName, templateType, jobDescription, userProfileSummary, mailSubject, mailBody, mailSubject, jobUrls)

	if err != nil {
		return "", "", err
	}

	return mailSubject, mailBody, nil
}

// TailorResumeWithJobDescriptionLLM generates a tailored, ATS-friendly resume in Markdown format.
func (co *Core) TailorResumeWithJobDescriptionLLM(jobDescription, extractedContent, companyName, jobRole string) (string, error) {
	ctx, cancel := getContextWithTimeout(30)
	defer cancel()

	res, err := co.llm.Models.GenerateContent(ctx, co.opts.ModelName, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			{Text: "[Job Description]:\n" + jobDescription + "\n[Company Name]:\n" + companyName + "\n[Job Role]:\n" + jobRole + "\n[Extracted Resume Content]:\n" + extractedContent},
			{Text: `
[Backstory]:
You are an expert FAANG resume strategist.

[Task]: Given a "Job Description", "Company Name", "Job Role", and "Extracted Resume Content", generate a concise, single-page, ATS-friendly Software Engineer resume in Markdown.

[Requirements]:
- Start with candidate's name as H1 and contact info (email, LinkedIn, GitHub, phone).
- The LinkedIn in the contact info MUST be a Markdown hyperlink with the full https URL, using the format: [complete.url](https://complete.url). Do NOT just write the URL or plain text; strictly use the Markdown hyperlink format.
- Add GitHub in the contact info MUST be a Markdown hyperlink with the full https URL, using the format: [github.com/sounishnath003](https://github.com/sounishnath003). Do NOT just write the URL or plain text; strictly use the Markdown hyperlink format.
- Email Address, LinkedIn, GitHub and Contact Number, must be in single Line. Separated by '|'.
- Add a brief "Professional Summary" tailored to the job description, company, and job role, using relevant keywords.
- List grouped skills (Languages, Frameworks, Cloud/DevOps, Tools) as bullet points.
- Show up to 3 most relevant roles (reverse-chronological), each with 3-5 quantified, action-oriented bullets (STAR/XYZ style).
- Select and include up to 2 personal projects that are most relevant to the job description and company. For each, paraphrase the project description and impact to closely align with the job requirements, company, and keywords. Present each project with a title and 2-3 concise, action-oriented bullet points, emphasizing technologies, outcomes, and relevance to the target role and company. Must Hyperlink Project Demo URL, if exists like [Github | YouTube | Link](https://complete.url).
- Focus only on content matching the job, company, and role; omit unrelated details.
- Use standard Markdown (no tables, no images, no extra commentary).
- Use present tense for current role, past tense for previous.
- Each bullet starts with a strong verb and includes metrics where possible.
- Max 1 page, ≤650 words, highly relevant for SWE roles at FAANG-level companies.
`},
			{Text: "\n\nNote: DO NOT USE Repeatative Action Verbs/words, Always Use unique action verbsin the Work experiences or Project sections."},
		},
	}}, nil)

	if err != nil {
		return "", fmt.Errorf("unable to generate tailored resume: %w", err)
	}

	if len(res.Candidates) == 0 {
		return "", errors.New("empty response from model")
	}

	return printResponse(res), nil
}
