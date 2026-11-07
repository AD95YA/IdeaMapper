package api

import (
	"encoding/json"
	"idea_app/db"
	"idea_app/models"
	"testing"
	"time"
)

func TestGenerateBackupBlob(t *testing.T) {
	if err := db.Initialize(":memory:"); err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer db.Close()

	// 1. Create an important idea and an unimportant idea
	impIdea := &models.Idea{
		ID:          models.NewID(),
		Title:       "Important Idea",
		Description: "Should be in backup",
		IsImportant: true,
	}
	regularIdea := &models.Idea{
		ID:          models.NewID(),
		Title:       "Regular Idea",
		Description: "Should NOT be in backup",
		IsImportant: false,
	}
	db.CreateIdea(impIdea)
	db.CreateIdea(regularIdea)

	// Add MetaData to the important idea
	meta := &models.MetaData{
		ID:          models.NewID(),
		ParentID:    impIdea.ID,
		ParentType:  "Idea",
		Level:       1,
		Title:       "Important Idea Sub-Thought",
		Description: "Sub description",
	}
	db.CreateMetaData(meta)

	// 2. Create an important memory
	impMem := &models.Memory{
		ID:          models.NewID(),
		Title:       "Important Memory",
		Description: "Saved forever",
		Date:        time.Now(),
		IsImportant: true,
	}
	db.CreateMemory(impMem)

	// 3. Generate backup blob
	blob, err := GenerateBackupBlob()
	if err != nil {
		t.Fatalf("GenerateBackupBlob failed: %v", err)
	}

	var payload BackupPayload
	if err := json.Unmarshal(blob, &payload); err != nil {
		t.Fatalf("Failed to parse JSON backup blob: %v", err)
	}

	if len(payload.Ideas) != 1 || payload.Ideas[0].Title != "Important Idea" {
		t.Errorf("Expected 1 important idea in backup, got: %d", len(payload.Ideas))
	}
	if len(payload.Ideas[0].MetaData) != 1 {
		t.Errorf("Expected 1 metadata attached to important idea, got: %d", len(payload.Ideas[0].MetaData))
	}
	if len(payload.Memories) != 1 || payload.Memories[0].Title != "Important Memory" {
		t.Errorf("Expected 1 important memory in backup, got: %d", len(payload.Memories))
	}
}
