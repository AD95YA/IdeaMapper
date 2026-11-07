package ai

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// GenerateIdeas calls the Gemini API to spawn new ideas based on a prompt using the environment API key
func GenerateIdeas(prompt string) (string, error) {
	return GenerateIdeasWithKey(prompt, "")
}

// GenerateIdeasWithKey calls the Gemini API using the provided key or falls back to GEMINI_API_KEY env
func GenerateIdeasWithKey(prompt string, apiKey string) (string, error) {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		return "", fmt.Errorf("gemini API key not configured")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text("Generate creative sub-thoughts or ideas for: "+prompt), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate ideas: %w", err)
	}

	text := resp.Text()
	if text == "" {
		return "", fmt.Errorf("no content generated")
	}
	return text, nil
}
