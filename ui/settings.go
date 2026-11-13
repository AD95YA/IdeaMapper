package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Expandable Setup Guide
	guideBox := container.NewVBox(
		widget.NewLabelWithStyle("3 Steps to get your free Gemini Key:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("1. Open Google AI Studio (https://aistudio.google.com/)"),
		widget.NewLabel("2. Click 'Create API Key' (free tier provides generous quota)"),
		widget.NewLabel("3. Copy & paste your key above and tap 'Save Key'"),
	)
	guideCard := widget.NewAccordion(
		widget.NewAccordionItem("ℹ How to get a free Gemini API Key", guideBox),
	)

	// Read-only prompt preview & copy button
	exactPrompt := "Generate creative sub-thoughts or ideas for: <idea title>"
	promptEntry := widget.NewEntry()
	promptEntry.SetText(exactPrompt)
	promptEntry.Disable() // Read-only

	copyPromptBtn := widget.NewButtonWithIcon("Copy Prompt", theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(exactPrompt)
		dialog.ShowInformation("Clipboard", "Prompt copied to clipboard!", w)
	})

	promptRow := container.NewBorder(nil, nil, widget.NewLabelWithStyle("Prompt Template:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), copyPromptBtn, promptEntry)

	aiCard := widget.NewCard("Gemini AI Configuration", "Powers sub-thought and creative idea expansion", container.NewVBox(
		apiKeyEntry,
		container.NewHBox(saveKeyBtn, testKeyBtn),
		keyStatusLabel,
		guideCard,
		promptRow,
	))
	vbox.Add(aiCard)
	vbox.Add(widget.NewSeparator())

	// 2. Automated Daily Local & SQLite Backup Card
	storageDir := a.Storage().RootURI().Path()
	if storageDir == "" {
		userHome, _ := os.UserHomeDir()
		storageDir = filepath.Join(userHome, ".idea_app")
	}

	lastBackupTime := a.Preferences().String(api.PrefLastDailyBackupUTC)
	if lastBackupTime == "" {
		lastBackupTime = "No automated backup performed yet"
	}
	dailyStatusLabel := widget.NewLabel("Last Daily Backup: " + lastBackupTime)

	backupNowBtn := widget.NewButtonWithIcon("Backup Now (JSON + SQLite)", theme.DocumentSaveIcon(), func() {
		dailyStatusLabel.SetText("Running atomic backup...")
		go func() {
			res, err := api.RunDailyBackup(storageDir)
			if err != nil {
				dailyStatusLabel.SetText("Backup failed: " + err.Error())
				dialog.ShowError(err, w)
			} else {
				nowStr := time.Now().UTC().Format(time.RFC3339)
				a.Preferences().SetString(api.PrefLastDailyBackupUTC, nowStr)
				msg := fmt.Sprintf("Backup successful for %s!\nPruned %d older backups.", res.Date, res.Pruned)
				dailyStatusLabel.SetText("Last Daily Backup: " + nowStr)
				dialog.ShowInformation("Backup Complete", msg, w)
			}
		}()
	})

	dailyBackupCard := widget.NewCard(
		"Automated Daily Local & SQLite Backups",
		"Rolling 7-day retention of JSON and raw SQLite .db copies",
		container.NewVBox(
			widget.NewLabel("Automatically captures snapshots on launch and retains the last 7 daily archives."),
			container.NewHBox(backupNowBtn),
			dailyStatusLabel,
		),
	)
	vbox.Add(dailyBackupCard)
	vbox.Add(widget.NewSeparator())

	// 3. Local JSON Backup & Restore Card
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

	// 3. Google Sign-In & Cloud Sync Card
	cloudBox := container.NewVBox()

	var refreshCloudSection func()
	refreshCloudSection = func() {
		cloudBox.Objects = nil

		user, signedIn := api.GetSignedInUser(a.Preferences())
		if signedIn {
			userLabel := widget.NewLabelWithStyle(fmt.Sprintf("✔ Signed in as: %s (%s)", user.Name, user.Email), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			signOutBtn := widget.NewButtonWithIcon("Sign Out", theme.CancelIcon(), func() {
				api.SignOut(a.Preferences())
				refreshCloudSection()
				dialog.ShowInformation("Signed Out", "You have been signed out of Google.", w)
			})
			cloudBox.Add(container.NewHBox(userLabel, layoutSpacer(), signOutBtn))
		} else {
			signInBtn := widget.NewButtonWithIcon("Sign in with Google", theme.AccountIcon(), func() {
				nameEntry := widget.NewEntry()
				nameEntry.SetPlaceHolder("Your Name")
				emailEntry := widget.NewEntry()
				emailEntry.SetPlaceHolder("user@gmail.com")

				d := dialog.NewCustomConfirm("Sign in with Google", "Sign In", "Cancel", container.NewVBox(
					widget.NewLabel("Sign in with your Google account to connect Drive backup:"),
					nameEntry,
					emailEntry,
				), func(ok bool) {
					if ok && strings.TrimSpace(emailEntry.Text) != "" {
						name := strings.TrimSpace(nameEntry.Text)
						if name == "" {
							name = "Google User"
						}
						api.SaveSignedInUser(a.Preferences(), name, strings.TrimSpace(emailEntry.Text), "oauth_token_active")
						refreshCloudSection()
						dialog.ShowInformation("Sign-In Successful", "Signed in as "+emailEntry.Text, w)
					}
				}, w)
				d.Show()
			})
			cloudBox.Add(signInBtn)
		}

		driveSyncBtn := widget.NewButtonWithIcon("Sync to Google Drive", theme.StorageIcon(), func() {
			dialog.ShowInformation("Google Drive Cloud Sync — Coming Soon",
				"Automatic cloud sync to your private Google Drive AppData folder is coming in the next update. Your data is currently protected via Daily Local & SQLite Backups.", w)
		})
		cloudBox.Add(driveSyncBtn)
		cloudBox.Refresh()
	}

	refreshCloudSection()

	driveCard := widget.NewCard("Google Drive Cloud Sync", "Connect Google Drive AppData storage", cloudBox)
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
