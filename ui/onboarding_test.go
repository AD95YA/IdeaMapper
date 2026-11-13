package ui

import (
	"testing"

	"idea_app/api"
)

func TestOnboarding_CorruptedJSONRejection(t *testing.T) {
	// Attempt importing invalid/corrupted JSON
	corrupted := []byte(`{"ideas": [{"title": "Broken", incomplete json...`)
	ideas, mems, err := api.ImportBackupBlob(corrupted)
	if err == nil {
		t.Fatalf("expected error importing corrupted JSON")
	}
	if ideas != 0 || mems != 0 {
		t.Fatalf("expected 0 imported items on corrupted payload, got %d, %d", ideas, mems)
	}
}

