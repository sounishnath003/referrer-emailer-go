package core

import (
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

	client, err := genai.NewClient(ctx, co.gcpProjectID, co.gcpLocation)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	co.llm = client.GenerativeModel(co.modelName)

	return nil
}
