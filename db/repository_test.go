package db

import (
	"idea_app/models"
	"strings"
	"testing"
	"time"
)

func TestMainDB(t *testing.T) {
	// 1. Setup in-memory DB
	err := Initialize(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer Close()

	// 2. Smoke Test: Create Idea
	idea := &models.Idea{
		ID:          models.NewID(),
		Title:       "Test Idea",
		Description: "Smoke Test",
	}
	if err := CreateIdea(idea); err != nil {
		t.Errorf("Failed to create idea: %v", err)
	}

	// 3. SQL Injection Test
	// Name contains malicious SQL. If prepared statements work, it will literally save this string.
	maliciousTitle := "Hack'; DROP TABLE ideas;--"
	maliciousIdea := &models.Idea{
		ID:    models.NewID(),
		Title: maliciousTitle,
	}
	if err := CreateIdea(maliciousIdea); err != nil {
		t.Errorf("SQL Injection test failed incorrectly: %v", err)
	}

	// Verify table still exists and data is safe
	ideas, err := GetAllIdeas()
	if err != nil || len(ideas) != 2 {
		t.Errorf("Ideas table was dropped or compromised by SQL Injection")
	}

	// 4. Test Memories CRUD
	mem := &models.Memory{
		ID:          models.NewID(),
		Title:       "Test Memory",
		Description: "Memory description",
		Date:        time.Now(),
		IsImportant: true,
	}
	if err := CreateMemory(mem); err != nil {
		t.Errorf("Failed to create memory: %v", err)
	}
	memories, err := GetAllMemories()
	if err != nil || len(memories) != 1 {
		t.Errorf("Failed to retrieve memories, got count: %d", len(memories))
	}

	// 5. Nested MetaData Logic & Enforcement (Corner Cases)
	n1 := &models.MetaData{ID: models.NewID(), ParentID: idea.ID, ParentType: "Idea", Level: 1, Title: "L1"}
	n2 := &models.MetaData{ID: models.NewID(), ParentID: n1.ID, ParentType: "MetaData", Level: 2, Title: "L2"}
	n3 := &models.MetaData{ID: models.NewID(), ParentID: n2.ID, ParentType: "MetaData", Level: 3, Title: "L3"}
	n4 := &models.MetaData{ID: models.NewID(), ParentID: n3.ID, ParentType: "MetaData", Level: 4, Title: "L4"}
	n5 := &models.MetaData{ID: models.NewID(), ParentID: n4.ID, ParentType: "MetaData", Level: 5, Title: "L5"}
	n6 := &models.MetaData{ID: models.NewID(), ParentID: n5.ID, ParentType: "MetaData", Level: 6, Title: "L6"}

	for _, n := range []*models.MetaData{n1, n2, n3, n4, n5, n6} {
		if err := CreateMetaData(n); err != nil {
			t.Errorf("Failed to insert valid MetaData level %d: %v", n.Level, err)
		}
	}

	// Negative Test: Exceeding max levels
	n7 := &models.MetaData{ID: models.NewID(), ParentID: n6.ID, ParentType: "MetaData", Level: 7, Title: "L7"}
	if err := CreateMetaData(n7); err == nil {
		t.Errorf("Allowed insertion of Level 7 MetaData!")
	} else if !strings.Contains(err.Error(), "maximum metadata nesting level") {
		t.Errorf("Unexpected error when testing Level 7 boundary: %v", err)
	}

	// Fetch MetaData hierarchy (Smoke test recursion)
	loadedMetaData, err := GetMetaDataForParent(idea.ID, "Idea")
	if err != nil || len(loadedMetaData) == 0 {
		t.Errorf("Failed to retrieve nested metaData: %v", err)
	}
	if len(loadedMetaData[0].SubMetaData) == 0 {
		t.Errorf("Expected recursive nested SubMetaData, but got none")
	}

	// 6. Test InterLinks
	link := &models.InterLink{
		ID:         models.NewID(),
		SourceID:   idea.ID,
		SourceType: "Idea",
		TargetID:   mem.ID,
		TargetType: "Memory",
		Name:       "Idea inspired memory",
	}
	if err := CreateInterLink(link); err != nil {
		t.Errorf("Failed to create interlink: %v", err)
	}
	linkMeta := &models.MetaData{
		ID:         models.NewID(),
		ParentID:   link.ID,
		ParentType: "InterLink",
		Level:      1,
		Title:      "Link note",
	}
	if err := CreateMetaData(linkMeta); err != nil {
		t.Errorf("Failed to create link metadata: %v", err)
	}
	links, err := GetInterLinksForSource(idea.ID)
	if err != nil || len(links) != 1 {
		t.Errorf("Failed to get interlinks for idea, count: %d", len(links))
	} else if links[0].LinkMetaData == nil {
		t.Errorf("LinkMetaData was not loaded for interlink")
	}

	// 7. Negative Test: NPEs and Nil data handling
	err = DeleteIdea("") // Deleting non-existent/empty ID
	if err != nil {
		t.Errorf("Empty deletion should safely perform NOOP, got error: %v", err)
	}

	// 8. Cascade Delete Test
	AddAttachment(n6.ID, "/test/path/img.png", true, models.NewID())

	// Trigger cascade
	if err := DeleteIdea(idea.ID); err != nil {
		t.Errorf("Failed to delete idea: %v", err)
	}

	// Verify cascade happened
	n6Deleted, err := GetMetaDataByID(n6.ID)
	if err != nil {
		t.Errorf("DB error querying deleted metadata: %v", err)
	}
	if n6Deleted != nil {
		t.Errorf("Cascade delete failed, submetadata L6 still exists")
	}

	// Verify attachments were deleted natively by Foreign Keys
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM attachments").Scan(&count)
	if count != 0 {
		t.Errorf("Cascade delete failed to drop attachments, count is %d", count)
	}
}
