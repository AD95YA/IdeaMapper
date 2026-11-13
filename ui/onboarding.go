package ui

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/api"
	"idea_app/db"
	"idea_app/models"
)

const PrefOnboardingCompleted = "onboarding_completed"

// CheckAndShowFirstLaunch prompts the user on first launch to restore from backup or start fresh
func CheckAndShowFirstLaunch(a fyne.App, w fyne.Window, onReady func()) {
	prefs := a.Preferences()
	if prefs.BoolWithFallback(PrefOnboardingCompleted, false) {
		if onReady != nil {
			onReady()
		}
		return
	}

	// Also check if existing records are present
	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()
	if len(ideas) > 0 || len(memories) > 0 {
		prefs.SetBool(PrefOnboardingCompleted, true)
		if onReady != nil {
			onReady()
		}
		return
	}

	var d dialog.Dialog

	importBtn := widget.NewButtonWithIcon("Import JSON Backup", theme.UploadIcon(), func() {
		d.Hide()
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				// Re-show onboarding prompt if cancelled
				CheckAndShowFirstLaunch(a, w, onReady)
				return
			}
			data, readErr := io.ReadAll(reader)
			reader.Close()
			if readErr != nil {
				dialog.ShowError(readErr, w)
				CheckAndShowFirstLaunch(a, w, onReady)
				return
			}

			ideasCount, memsCount, importErr := api.ImportBackupBlob(data)
			if importErr != nil {
				dialog.ShowError(fmt.Errorf("invalid backup file: %w", importErr), w)
				CheckAndShowFirstLaunch(a, w, onReady)
				return
			}

			prefs.SetBool(PrefOnboardingCompleted, true)
			msg := fmt.Sprintf("Restored %d Ideas and %d Memories successfully!", ideasCount, memsCount)
			dialog.ShowInformation("Restore Complete", msg, w)
			if onReady != nil {
				onReady()
			}
		}, w)
		fd.Show()
	})
	importBtn.Importance = widget.HighImportance

	startFreshBtn := widget.NewButtonWithIcon("Start Fresh", theme.ContentAddIcon(), func() {
		d.Hide()
		// Initialize starter root topic
		starter := &models.Idea{
			ID:          models.NewID(),
			Title:       "Welcome to Mendu: IdeaMapper! 🚀",
			Description: "Your personal, offline-first knowledge graph. Tap 'Add Sub-Thought' to branch ideas into multi-level orbits, or connect ideas and memories in the Mind Map!",
			IsImportant: true,
		}
		_ = db.CreateIdea(starter)
		prefs.SetBool(PrefOnboardingCompleted, true)
		if onReady != nil {
			onReady()
		}
	})

	content := widget.NewCard(
		"Welcome to Mendu: IdeaMapper",
		"Personal Knowledge Graph & Memory Journal",
		widget.NewLabel("Would you like to restore an existing JSON backup or start a fresh mind map?"),
	)

	d = dialog.NewCustomWithoutButtons("First Launch", content, w)
	content.SetContent(widget.NewCard(
		"Welcome to Mendu: IdeaMapper",
		"Personal Knowledge Graph & Memory Journal",
		container.NewVBox(
			widget.NewLabel("Would you like to restore an existing JSON backup or start a fresh mind map?"),
			widget.NewSeparator(),
			container.NewHBox(importBtn, layoutSpacer(), startFreshBtn),
		),
	))
	d.Show()
}
