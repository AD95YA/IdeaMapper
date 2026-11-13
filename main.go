package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"idea_app/api"
	"idea_app/db"
	"idea_app/ui"
)

func main() {
	a := app.NewWithID("com.antigravity.ideamemoryapp")
	w := a.NewWindow("Mendu: IdeaMapper")

	// Set up database in app storage
	storageDir := a.Storage().RootURI().Path() // Platform-safe storage directory
	if storageDir == "" {
		// Fallback for simple local testing if Fyne URI isn't file-based
		userDir, _ := os.UserHomeDir()
		storageDir = filepath.Join(userDir, ".idea_app")
		os.MkdirAll(storageDir, 0755)
	}

	if err := db.Initialize(storageDir); err != nil {
		dialog.ShowError(err, w)
		log.Fatalf("Database Init Error: %v", err)
	}
	defer db.Close()

	// 1. Start Background Daily Backup Scheduler (JSON + SQLite, 7-day rolling)
	api.StartDailyBackupScheduler(a, storageDir)

	// 2. Load Initial UI Layout
	w.SetContent(ui.BuildMainTabs(a, w))
	w.Resize(fyne.NewSize(400, 700)) // Mobile-like aspect ratio for desktop testing

	// 3. First-Launch Onboarding Check
	ui.CheckAndShowFirstLaunch(a, w, func() {
		w.SetContent(ui.BuildMainTabs(a, w))
	})

	w.ShowAndRun()
}
