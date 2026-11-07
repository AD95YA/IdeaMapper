package api

import (
	"encoding/json"
	"idea_app/db"
	"idea_app/models"
	"testing"
	"time"
)

func TestGenerateAndImportBackupBlob(t *testing.T) {
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

	// 4. Test Import into a fresh clean database
	db.Close()
	if err := db.Initialize(":memory:"); err != nil {
		t.Fatalf("Failed to re-init DB: %v", err)
	}

	ideasImported, memoriesImported, err := ImportBackupBlob(blob)
	if err != nil {
		t.Fatalf("ImportBackupBlob failed: %v", err)
	}
	if ideasImported != 1 || memoriesImported != 1 {
		t.Errorf("Expected 1 idea and 1 memory imported, got %d ideas, %d memories", ideasImported, memoriesImported)
	}

	// Verify imported idea has its metadata restored
	restoredIdea, err := db.GetIdeaByID(impIdea.ID)
	if err != nil || restoredIdea == nil {
		t.Fatalf("Restored idea not found")
	}
	if len(restoredIdea.MetaData) != 1 || restoredIdea.MetaData[0].Title != "Important Idea Sub-Thought" {
		t.Errorf("Restored metadata missing or incorrect")
	}
}
