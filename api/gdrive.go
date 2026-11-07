package api

import (
	"encoding/json"
	"fmt"
	"idea_app/db"
	"idea_app/models"
)

type BackupPayload struct {
	Ideas    []*models.Idea   `json:"ideas"`
	Memories []*models.Memory `json:"memories"`
}

// GenerateBackupBlob extracts all IsImportant=true entities and translates them to JSON
func GenerateBackupBlob() ([]byte, error) {
	ideas, err := db.GetAllIdeas()
	if err != nil {
		return nil, err
	}

	memories, err := db.GetAllMemories()
	if err != nil {
		return nil, err
	}

	var impIdeas []*models.Idea
	for _, i := range ideas {
		if i.IsImportant {
			// Fetch its metaData too
			metaData, _ := db.GetMetaDataForParent(i.ID, "Idea")
			i.MetaData = metaData
			impIdeas = append(impIdeas, i)
		}
	}

	var impMemories []*models.Memory
	for _, m := range memories {
		if m.IsImportant {
			metaData, _ := db.GetMetaDataForParent(m.ID, "Memory")
			m.MetaData = metaData
			impMemories = append(impMemories, m)
		}
	}

	payload := BackupPayload{Ideas: impIdeas, Memories: impMemories}

	return json.MarshalIndent(payload, "", "  ")
}

// ImportBackupBlob restores entities from JSON into the database
func ImportBackupBlob(jsonData []byte) (int, int, error) {
	var payload BackupPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return 0, 0, fmt.Errorf("invalid backup format: %w", err)
	}

	ideasImported := 0
	for _, idea := range payload.Ideas {
		existing, _ := db.GetIdeaByID(idea.ID)
		if existing == nil {
			if err := db.CreateIdea(idea); err == nil {
				ideasImported++
			}
		}
		restoreMetaData(idea.MetaData, idea.ID, "Idea")
	}

	memoriesImported := 0
	for _, mem := range payload.Memories {
		existing, _ := db.GetMemoryByID(mem.ID)
		if existing == nil {
			if err := db.CreateMemory(mem); err == nil {
				memoriesImported++
			}
		}
		restoreMetaData(mem.MetaData, mem.ID, "Memory")
	}

	return ideasImported, memoriesImported, nil
}

func restoreMetaData(metaList []*models.MetaData, parentID string, parentType string) {
	for _, m := range metaList {
		m.ParentID = parentID
		m.ParentType = parentType
		existing, _ := db.GetMetaDataByID(m.ID)
		if existing == nil {
			_ = db.CreateMetaData(m)
		}
		if len(m.SubMetaData) > 0 {
			restoreMetaData(m.SubMetaData, m.ID, "MetaData")
		}
	}
}

// UploadToDrive is a placeholder for the actual OAuth2 drive.v3 API upload.
func UploadToDrive(jsonData []byte, authCode string) error {
	// 1. Exchange authCode for Token
	// 2. Initialize drive.NewService
	// 3. Create drive.File metadata
	// 4. Call driveService.Files.Create(file).Media(bytes.NewReader(jsonData)).Do()
	return fmt.Errorf("Google Drive API OAuth integration requires specific client configurations")
}
