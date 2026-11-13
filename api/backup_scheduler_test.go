package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"idea_app/db"
	"idea_app/models"
)

func TestShouldRunDailyBackup(t *testing.T) {
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

	// 1. Never backed up -> Should run
	if !ShouldRunDailyBackup("", now) {
		t.Fatalf("expected backup to run when lastBackupISO is empty")
	}

	// 2. Backed up earlier today (2 hours ago) -> Should NOT run
	earlierToday := time.Date(2026, 11, 13, 10, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if ShouldRunDailyBackup(earlierToday, now) {
		t.Fatalf("expected backup to skip when already run today")
	}

	// 3. Backed up yesterday (26 hours ago) -> Should run
	yesterday := time.Date(2026, 11, 12, 10, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if !ShouldRunDailyBackup(yesterday, now) {
		t.Fatalf("expected backup to run when >= 24 hours elapsed")
	}

	// 4. Invalid date format -> Should run as fallback
	if !ShouldRunDailyBackup("corrupted-date-string", now) {
		t.Fatalf("expected backup to run on invalid date format")
	}
}

func TestPruneRollingBackups(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock backup files for 10 consecutive days
	baseDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		dateStr := baseDate.AddDate(0, 0, i).Format("2006-01-02")
		jsonFile := filepath.Join(tempDir, fmt.Sprintf("ideamapper_%s.json", dateStr))
		dbFile := filepath.Join(tempDir, fmt.Sprintf("ideamapper_%s.db", dateStr))

		_ = os.WriteFile(jsonFile, []byte("json"), 0644)
		_ = os.WriteFile(dbFile, []byte("sqlite"), 0644)
	}

	// Keep only the 7 newest days
	pruned, err := PruneRollingBackups(tempDir, 7)
	if err != nil {
		t.Fatalf("prune failed: %v", err)
	}

	// 3 days * 2 files each = 6 files pruned
	if pruned != 6 {
		t.Fatalf("expected 6 files pruned, got %d", pruned)
	}

	// The 3 oldest days (Nov 01, 02, 03) should be deleted
	for i := 0; i < 3; i++ {
		dateStr := baseDate.AddDate(0, 0, i).Format("2006-01-02")
		jsonFile := filepath.Join(tempDir, fmt.Sprintf("ideamapper_%s.json", dateStr))
		if _, err := os.Stat(jsonFile); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be deleted", jsonFile)
		}
	}

	// The 7 newest days (Nov 04 through 10) should remain
	for i := 3; i < 10; i++ {
		dateStr := baseDate.AddDate(0, 0, i).Format("2006-01-02")
		jsonFile := filepath.Join(tempDir, fmt.Sprintf("ideamapper_%s.json", dateStr))
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			t.Fatalf("expected %s to be preserved", jsonFile)
		}
	}
}

func TestRunDailyBackup_FullRoundtrip(t *testing.T) {
	tempDir := t.TempDir()

	// Initialize test database
	err := db.Initialize(tempDir)
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}
	defer db.Close()

	// Seed one idea
	idea := &models.Idea{
		ID:          "test-idea-backup",
		Title:       "Backup Test Idea",
		Description: "Testing daily backup serialization",
		IsImportant: false, // Should be included in FullBackupBlob
	}
	_ = db.CreateIdea(idea)

	res, err := RunDailyBackup(tempDir)
	if err != nil {
		t.Fatalf("RunDailyBackup failed: %v", err)
	}

	if res == nil {
		t.Fatalf("expected non-nil BackupResult")
	}

	// Verify JSON file exists and contains the idea
	data, err := os.ReadFile(res.JSONPath)
	if err != nil {
		t.Fatalf("failed to read generated JSON backup: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty JSON backup content")
	}

	// Verify DB file exists
	dbInfo, err := os.Stat(res.DBPath)
	if err != nil {
		t.Fatalf("failed to stat generated DB backup: %v", err)
	}
	if dbInfo.Size() == 0 {
		t.Fatalf("expected non-empty DB backup file")
	}
}
