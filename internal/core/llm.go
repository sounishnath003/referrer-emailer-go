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

	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  resumePath,
	}

	res, err := co.llm.GenerateContent(ctx, part, genai.Text(`
		[Story]: You are a very professional Hiring Manager specialist of a Company. You have to screen candidate based on their resume.
		
		[Task]: Given a "Resume", You have to extract all the "Key Informations" of this candidate.

		[Important to note]:
    	1. Extract all the Social Links (Linkedin, Github, MobileNo, EmailAdress etc).
		2. Extract the "Skillsets" which are important for a Job description to match.
		3. Extract all "Work Experiences" of this candidate in Chronological Order. Along with the "Key Performance" highlights.
		4. Extract "Achievements" of this candidates if present. Else "SKIP".
		5. Extract "Latest Education" (College, Institutions), Year of Graduation Tenure. Skip the Class "10th, 12th" portions.
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

// GenerateProfileSummary generates a summary of a resume using an AI model.
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
func (co *Core) GenerateProfileSummary(resumePath string) error {
	co.Lo.Info("started parsing", "resume", resumePath)

	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  resumePath,
	}

	res, err := co.llm.GenerateContent(ctx, part, genai.Text(`
			You are a very professional document summarization specialist.
    		Please summarize the given document.
	`))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	co.Lo.Info("generated by AI", "response", res.Candidates[0].Content.Parts[0])
	return nil
}

// ResumeParser parses the resume from the given file path and returns the plain text content.
func (co *Core) ResumeParser(resumePath string) error {
	co.Lo.Info("started parsing", "resume", resumePath)

	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  "gs://sounish-cloud-workstation/referrer-uploads/Sounish_Nath_Resume_25.pdf",
	}

	res, err := co.llm.GenerateContent(ctx, part, genai.Text(`
			You are a very professional document summarization specialist.
    		Please summarize the given document.
	`))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	co.Lo.Info("generated by AI", "response", res.Candidates[0].Content.Parts[0])
	return nil
}
