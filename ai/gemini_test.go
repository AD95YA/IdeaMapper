package ai

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateIdeasWithKey_Validation(t *testing.T) {
	// 1. Empty key and no environment key -> Should return descriptive error
	_ = os.Unsetenv("GEMINI_API_KEY")
	_, err := GenerateIdeasWithKey("Test Prompt", "")
	if err == nil {
		t.Fatalf("expected error when GEMINI_API_KEY is not configured")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("unexpected error message: %v", err)
	}

	// 2. GenerateIdeas wrapper with no key configured
	_, err = GenerateIdeas("Another Prompt")
	if err == nil {
		t.Fatalf("expected error from GenerateIdeas without key")
	}
}
