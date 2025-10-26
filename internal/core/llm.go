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
		Project: co.opts.GcpProjectID, Location: co.opts.GcpLocation, Backend: genai.BackendVertexAI},
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
				result += fmt.Sprint(part.Text)
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
		[Role]: You are an expert HR Tech Specialist responsible for parsing and structuring resume data.

		[Task]:
		Extract all key information from the provided resume and structure it for a database.

		[Instructions]:
		1.  **Candidate Name**: Start with "## Candidate Name". Do not add any text before it.
		2.  **Contact Information**: Extract email, phone number, and personal website/portfolio.
		3.  **Social Links**: Extract all social media links (e.g., LinkedIn, GitHub, Twitter).
		4.  **Skillsets**:
			-   Categorize skills into: 'Programming Languages', 'Frameworks & Libraries', 'Databases', 'Cloud & DevOps', and 'Tools'.
			-   List skills under their respective categories.
		5.  **Work Experience**:
			-   Detail all work experiences in reverse chronological order.
			-   For each role, include: 'Job Title', 'Company', 'Location', 'Dates of Employment', and 3-5 bullet points describing key responsibilities and quantifiable achievements (e.g., "Increased API response time by 30%").
		6.  **Projects**:
			-   List all personal or professional projects.
			-   For each project, include: 'Project Name', 'Technologies Used', and links to demos or source code (e.g., GitHub, live URL).
		7.  **Achievements**: If any awards or recognitions are mentioned, extract them. If not, omit this section.
		8.  **Education**:
			-   Extract the most recent educational qualifications.
			-   Include: 'University', 'Degree', 'Field of Study', and 'Graduation Year'.
			-   Omit high school details.

		[Output Format]:
		The output must be a well-structured markdown document.
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
			[Role]: You are a Professional Content Generation Specialist.

			[Task]: Summarize the given resume content into concise bullet points.

			[Instructions]:
			1.  Focus on professional work experience, skills, projects, and achievements.
			2.  Ensure contact details (phone, email, LinkedIn, portfolio, etc.) are included in the summary.
			3.  The entire output must be in Markdown format.
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
		[Role]: You are a Senior Data Entry Specialist.

		[Task]: Convert the provided resume content into a structured JSON format.

		[Instructions]:
		1.  Adhere strictly to the JSON schema provided below.
		2.  The JSON output must be a single line of text without any special characters or formatting.
		3.  Do not include any explanatory text or markdown formatting in your response.

		[JSON Schema]:
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
	ctx, cancel := getContextWithTimeout(60)
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
			[Role]: You are a professional career coach and expert cold-email copywriter drafting an email on behalf of a candidate.

			[Task]: Write a compelling cold email to a recruiter to secure an interview for a job opportunity.

			[Instructions]:
			1.  **Tone and Style**: Write in the first person from the candidate's perspective. The tone must be professional, concise, and genuinely enthusiastic.
			2.  **Structure**:
				- **Opening**: Start with a strong opening that grabs the recruiter's attention.
				- **Body**: Highlight the candidate's most relevant skills and experiences, directly aligning them with the job description. Use the STAR method (Situation, Task, Action, Result) to frame 1-2 key achievements (e.g., "Achieved X by doing Y, resulting in Z").
				- **Call to Action**: End with a clear and confident call to action, suggesting a brief chat.
				- **Signature**: Include a professional signature with the candidate's full contact details (phone, email, LinkedIn, portfolio).
			3.  **Formatting**:
				- Use Markdown for clear formatting.
				- Emphasize key skills and achievements with bold keywords.
				- List the "JOB URLs" as a bulleted list if applicable.
			4.  **Constraints**:
				- The email body must be under 200 words to ensure it gets read.
				- Do not include the subject line in the output.
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
		[Role]: You are an expert copywriter specializing in email marketing.

		[Task]: Generate a concise, professional, and attention-grabbing email subject line based on the provided email body and job details.

		[Instructions]:
		1.  The subject line must be tailored to the email content, job description, and company.
		2.  It must explicitly mention the company name: "%s".
		3.  The output should ONLY be the subject line itself, without any extra text, quotes, or formatting.
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
[Role]: You are an expert FAANG resume strategist and ATS optimization specialist.

[Task]: Generate a concise, single-page, and strictly ATS-friendly Software Engineer resume in Markdown. The resume must be tailored to the provided job description, company, and role.

[Instructions]:
1.  **ATS-Friendliness is Priority**: Use standard formatting. Avoid tables, columns, and images. Use standard bullet points.
2.  **Header**:
    - Start with the candidate's name as an H1 heading.
    - Follow with contact information (Email, LinkedIn, GitHub, Phone) on a single line, separated by '|'.
    - LinkedIn and GitHub links MUST be in Markdown hyperlink format (e.g., [LinkedIn](https://linkedin.com/in/user)).
3.  **Professional Summary**:
    - Write a brief, impactful summary (2-3 sentences) tailored to the job, highlighting key skills and years of experience.
4.  **Skills**:
    - Group skills into logical categories (e.g., Languages, Frameworks, Cloud/DevOps, Tools).
    - List skills as bullet points, matching keywords from the job description where appropriate.
5.  **Work Experience**:
    - List up to 3 of the most relevant roles in reverse-chronological order.
    - For each role, provide 3-5 bullet points using the STAR or XYZ method to quantify achievements (e.g., "Increased performance by 20% by implementing X").
    - Use strong, unique action verbs and metrics. Use present tense for the current role and past tense for previous roles.
6.  **Projects**:
    - Include up to 2 of the most relevant personal projects.
    - For each project, provide a title and 2-3 bullet points describing the project, technologies used, and its relevance to the job.
    - Hyperlink project demo URLs where available (e.g., [GitHub](https://github.com/user/project)).

[Output Format]:
- The entire output must be in standard Markdown.
- Do not include any commentary or explanations outside of the resume content.

[Constraints]:
- **Strictly One Page**: The resume must not exceed one page (approximately 650 words).
- **No Exaggeration**: Present skills and achievements accurately and honestly. Do not invent or overstate qualifications.
- **Relevance is Key**: Omit any information not directly relevant to the target role.
`},
			{Text: "\n\nNote: Do not use repetitive action verbs. Ensure variety in the language used in the work experience and project sections."},
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
