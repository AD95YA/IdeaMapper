package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

const (
	PrefLastDailyBackupUTC = "last_daily_backup_utc"
	DefaultRetentionDays   = 7
)

// BackupResult encapsulates metadata about a completed backup run
type BackupResult struct {
	JSONPath string
	DBPath   string
	Date     string
	Pruned   int
}

// RunDailyBackup performs an atomic daily backup of JSON and SQLite database files
func RunDailyBackup(storageDir string) (*BackupResult, error) {
	if storageDir == "" {
		return nil, fmt.Errorf("storage directory cannot be empty")
	}

	backupDir := filepath.Join(storageDir, "backups", "daily")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	today := time.Now().UTC().Format("2006-01-02")
	targetJSON := filepath.Join(backupDir, fmt.Sprintf("ideamapper_%s.json", today))
	targetDB := filepath.Join(backupDir, fmt.Sprintf("ideamapper_%s.db", today))

	// 1. Generate & Write Full JSON Backup Atomically
	jsonBlob, err := GenerateFullBackupBlob()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JSON backup: %w", err)
	}

	tmpJSON := filepath.Join(backupDir, fmt.Sprintf(".tmp_%s.json", today))
	if err := os.WriteFile(tmpJSON, jsonBlob, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp JSON backup: %w", err)
	}
	if err := os.Rename(tmpJSON, targetJSON); err != nil {
		os.Remove(tmpJSON)
		return nil, fmt.Errorf("failed to finalize JSON backup file: %w", err)
	}

	// 2. Clone SQLite Database Atomically
	srcDB := filepath.Join(storageDir, "app.db")
	if _, err := os.Stat(srcDB); err == nil {
		tmpDB := filepath.Join(backupDir, fmt.Sprintf(".tmp_%s.db", today))
		if err := copyFile(srcDB, tmpDB); err != nil {
			return nil, fmt.Errorf("failed to copy SQLite database: %w", err)
		}
		if err := os.Rename(tmpDB, targetDB); err != nil {
			os.Remove(tmpDB)
			return nil, fmt.Errorf("failed to finalize SQLite backup file: %w", err)
		}
	}

	// 3. Prune Rolling Backups to Retain Last 7 Days
	prunedCount, err := PruneRollingBackups(backupDir, DefaultRetentionDays)
	if err != nil {
		// Log error but don't fail backup
		fmt.Printf("Warning: error during backup pruning: %v\n", err)
	}

	return &BackupResult{
		JSONPath: targetJSON,
		DBPath:   targetDB,
		Date:     today,
		Pruned:   prunedCount,
	}, nil
}

// ShouldRunDailyBackup checks if 24 hours have elapsed since the last backup
func ShouldRunDailyBackup(lastBackupISO string, now time.Time) bool {
	if lastBackupISO == "" {
		return true
	}
	lastTime, err := time.Parse(time.RFC3339, lastBackupISO)
	if err != nil {
		// Try parsing simple YYYY-MM-DD
		lastDate, err2 := time.Parse("2006-01-02", lastBackupISO)
		if err2 != nil {
			return true
		}
		lastTime = lastDate
	}

	// Check if >= 24 hours or different calendar day in UTC
	if now.UTC().Sub(lastTime.UTC()) >= 24*time.Hour {
		return true
	}
	if now.UTC().Format("2006-01-02") != lastTime.UTC().Format("2006-01-02") {
		return true
	}

	return false
}

// StartDailyBackupScheduler runs background backup checks on launch and resume
func StartDailyBackupScheduler(a fyne.App, storageDir string) {
	go func() {
		prefs := a.Preferences()
		lastBackup := prefs.String(PrefLastDailyBackupUTC)
		now := time.Now().UTC()

		if ShouldRunDailyBackup(lastBackup, now) {
			result, err := RunDailyBackup(storageDir)
			if err == nil && result != nil {
				prefs.SetString(PrefLastDailyBackupUTC, now.Format(time.RFC3339))
			}
		}
	}()
}

// PruneRollingBackups keeps the newest keepDays backup dates and deletes older files
func PruneRollingBackups(backupDir string, keepDays int) (int, error) {
	if keepDays <= 0 {
		keepDays = DefaultRetentionDays
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return 0, err
	}

	// Group files by extracted date prefix (ideamapper_YYYY-MM-DD.*)
	dateMap := make(map[string][]string)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "ideamapper_") && (strings.HasSuffix(name, ".json") || strings.HasSuffix(name, ".db")) {
			// Extract date part
			parts := strings.Split(name, "_")
			if len(parts) >= 2 {
				datePart := strings.Split(parts[1], ".")[0]
				dateMap[datePart] = append(dateMap[datePart], filepath.Join(backupDir, name))
			}
		}
	}

	var uniqueDates []string
	for d := range dateMap {
		uniqueDates = append(uniqueDates, d)
	}
	sort.Strings(uniqueDates)

	prunedCount := 0
	if len(uniqueDates) > keepDays {
		toDeleteDates := uniqueDates[:len(uniqueDates)-keepDays]
		for _, d := range toDeleteDates {
			for _, filePath := range dateMap[d] {
				if err := os.Remove(filePath); err == nil {
					prunedCount++
				}
			}
		}
	}

	return prunedCount, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

