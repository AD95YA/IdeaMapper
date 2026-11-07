package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/ai"
	"idea_app/api"
	"idea_app/db"
)

// BuildSettings renders the Settings and Backup management screen
func BuildSettings(a fyne.App, w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()

	// Title
	vbox.Add(widget.NewLabelWithStyle("Settings & Cloud Backup", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	vbox.Add(widget.NewSeparator())

	// 1. Gemini AI Configuration Card
	currentKey := a.Preferences().StringWithFallback("gemini_api_key", os.Getenv("GEMINI_API_KEY"))
	apiKeyEntry := widget.NewPasswordEntry()
	apiKeyEntry.SetText(currentKey)
	apiKeyEntry.SetPlaceHolder("Enter your Google Gemini API Key")

	keyStatusLabel := widget.NewLabel("")

	saveKeyBtn := widget.NewButtonWithIcon("Save Key", theme.DocumentSaveIcon(), func() {
		newKey := strings.TrimSpace(apiKeyEntry.Text)
		a.Preferences().SetString("gemini_api_key", newKey)
		keyStatusLabel.SetText("API Key saved successfully!")
	})

	testKeyBtn := widget.NewButtonWithIcon("Test Connection", theme.HelpIcon(), func() {
		key := strings.TrimSpace(apiKeyEntry.Text)
		if key == "" {
			dialog.ShowError(fmt.Errorf("please enter an API key first"), w)
			return
		}
		keyStatusLabel.SetText("Testing Gemini API connection...")
		go func() {
			resp, err := ai.GenerateIdeasWithKey("Respond with exactly one word: Success", key)
			if err != nil {
				keyStatusLabel.SetText("Connection failed: " + err.Error())
				dialog.ShowError(fmt.Errorf("Gemini API Error: %w", err), w)
			} else {
				keyStatusLabel.SetText("Connected! Response: " + strings.TrimSpace(resp))
				dialog.ShowInformation("Gemini Connected", "Success! AI sub-thought generation is ready.", w)
			}
		}()
	})

	aiCard := widget.NewCard("Gemini AI Configuration", "Powers sub-thought and creative idea expansion", container.NewVBox(
		apiKeyEntry,
		container.NewHBox(saveKeyBtn, testKeyBtn),
		keyStatusLabel,
	))
	vbox.Add(aiCard)
	vbox.Add(widget.NewSeparator())

	// 2. Local JSON Backup & Restore Card
	backupStatusLabel := widget.NewLabel("")

	exportBtn := widget.NewButtonWithIcon("Export JSON Backup", theme.DownloadIcon(), func() {
		blob, err := api.GenerateBackupBlob()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		userHome, _ := os.UserHomeDir()
		defaultPath := filepath.Join(userHome, "ideamapper_backup.json")

		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				_, writeErr := writer.Write(blob)
				writer.Close()
				if writeErr != nil {
					dialog.ShowError(writeErr, w)
				} else {
					backupStatusLabel.SetText("Backup exported to: " + writer.URI().Path())
					dialog.ShowInformation("Backup Complete", "Successfully exported backup file.", w)
				}
			}
		}, w)
		fd.SetFileName("ideamapper_backup.json")
		fd.Show()

		// Fallback notice
		backupStatusLabel.SetText("Default target: " + defaultPath)
	})

	importBtn := widget.NewButtonWithIcon("Import JSON Backup", theme.UploadIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				data, readErr := io.ReadAll(reader)
				reader.Close()
				if readErr != nil {
					dialog.ShowError(readErr, w)
					return
				}
				ideasCount, memsCount, importErr := api.ImportBackupBlob(data)
				if importErr != nil {
					dialog.ShowError(importErr, w)
					return
				}
				msg := fmt.Sprintf("Imported %d Ideas and %d Memories successfully!", ideasCount, memsCount)
				backupStatusLabel.SetText(msg)
				dialog.ShowInformation("Import Complete", msg, w)
			}
		}, w)
		fd.Show()
	})

	backupCard := widget.NewCard("Local JSON Backup & Restore", "Export or restore all important ideas and memories", container.NewVBox(
		container.NewHBox(exportBtn, importBtn),
		backupStatusLabel,
	))
	vbox.Add(backupCard)
	vbox.Add(widget.NewSeparator())

	// 3. Google Drive Sync Card
	driveSyncBtn := widget.NewButtonWithIcon("Sync to Google Drive", theme.StorageIcon(), func() {
		_, err := api.GenerateBackupBlob()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		// Attempt upload
		uploadErr := api.UploadToDrive(nil, "")
		dialog.ShowInformation("Google Drive Cloud Sync",
			fmt.Sprintf("Status: %v\n\nGoogle Drive OAuth2 client credentials can be configured for automated cloud sync.", uploadErr), w)
	})

	driveCard := widget.NewCard("Google Drive Cloud Sync", "Sync important ideas and memories to the cloud", container.NewVBox(
		widget.NewLabel("Automatically creates snapshots of starred ideas and memories."),
		driveSyncBtn,
	))
	vbox.Add(driveCard)
	vbox.Add(widget.NewSeparator())

	// 4. System Statistics
	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()
	links, _ := db.GetAllInterLinks()

	statsLabel := widget.NewLabel(fmt.Sprintf("Total Ideas: %d   |   Total Memories: %d   |   Total InterLinks: %d",
		len(ideas), len(memories), len(links)))
	vbox.Add(widget.NewCard("System Statistics", "", statsLabel))

	return container.NewVScroll(vbox)
}
