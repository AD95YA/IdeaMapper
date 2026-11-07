package ui

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/ai"
	"idea_app/db"
	"idea_app/models"
)

// ShowIdeaDetail renders the complete detail view for an Idea
func ShowIdeaDetail(a fyne.App, w fyne.Window, ideaID string, onBack func()) {
	var renderView func()

	renderView = func() {
		idea, err := db.GetIdeaByID(ideaID)
		if err != nil || idea == nil {
			dialog.ShowError(fmt.Errorf("idea not found or failed to load"), w)
			if onBack != nil {
				onBack()
			}
			return
		}

		// Top Bar
		backBtn := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
			if onBack != nil {
				onBack()
			}
		})

		var starIcon fyne.Resource
		starText := "Regular"
		if idea.IsImportant {
			starIcon = theme.ConfirmIcon()
			starText = "★ Important"
		} else {
			starIcon = theme.CancelIcon()
			starText = "☆ Normal"
		}

		importantBtn := widget.NewButtonWithIcon(starText, starIcon, func() {
			idea.IsImportant = !idea.IsImportant
			if err := db.UpdateIdea(idea); err != nil {
				dialog.ShowError(err, w)
				return
			}
			renderView()
		})

		editBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
			showEditIdeaModal(w, idea, renderView)
		})

		deleteBtn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Delete Idea", "Are you sure you want to delete '"+idea.Title+"'? All sub-thoughts will also be deleted.", func(ok bool) {
				if ok {
					db.DeleteIdea(idea.ID)
					if onBack != nil {
						onBack()
					}
				}
			}, w)
		})

		aiBtn := widget.NewButtonWithIcon("🤖 Ask Gemini AI", theme.HelpIcon(), func() {
			showAIGenerateModal(a, w, idea, renderView)
		})

		actionsBox := container.NewHBox(importantBtn, editBtn, aiBtn, deleteBtn)
		topBar := container.NewBorder(nil, nil, backBtn, actionsBox, widget.NewLabelWithStyle(idea.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

		// Body Content
		body := container.NewVBox()

		// Title & Description card
		descText := idea.Description
		if descText == "" {
			descText = "(No description provided)"
		}
		descCard := widget.NewCard("Description", "", widget.NewLabel(descText))
		body.Add(descCard)
		body.Add(widget.NewSeparator())

		// InterLinks Section
		linksSection := BuildInterLinksSection(a, w, idea.ID, "Idea", renderView, func(targetType string, targetID string) {
			if targetType == "Idea" {
				ShowIdeaDetail(a, w, targetID, renderView)
			} else {
				ShowMemoryDetail(a, w, targetID, renderView)
			}
		})
		body.Add(linksSection)
		body.Add(widget.NewSeparator())

		// MetaData Hierarchical Tree
		metaDataSection := BuildMetaDataTree(w, idea.ID, "Idea", renderView)
		body.Add(metaDataSection)

		scroll := container.NewVScroll(body)
		mainLayout := container.NewBorder(topBar, nil, nil, nil, scroll)

		w.SetContent(mainLayout)
	}

	renderView()
}

func showEditIdeaModal(w fyne.Window, idea *models.Idea, onSaved func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetText(idea.Title)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(idea.Description)

	iconEntry := widget.NewEntry()
	iconEntry.SetText(idea.IconPath)
	iconEntry.SetPlaceHolder("Optional icon identifier or path")

	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
		widget.NewFormItem("Icon Path", iconEntry),
	)

	d := dialog.NewCustomConfirm("Edit Idea", "Save", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			idea.Title = strings.TrimSpace(titleEntry.Text)
			idea.Description = strings.TrimSpace(descEntry.Text)
			idea.IconPath = strings.TrimSpace(iconEntry.Text)
			if err := db.UpdateIdea(idea); err != nil {
				dialog.ShowError(err, w)
				return
			}
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 280))
	d.Show()
}

func showAIGenerateModal(a fyne.App, w fyne.Window, idea *models.Idea, onSaved func()) {
	promptEntry := widget.NewMultiLineEntry()
	promptEntry.SetText(fmt.Sprintf("%s\n%s", idea.Title, idea.Description))

	apiKey := a.Preferences().StringWithFallback("gemini_api_key", os.Getenv("GEMINI_API_KEY"))

	resultBox := widget.NewMultiLineEntry()
	resultBox.SetPlaceHolder("Gemini suggestions will appear here...")
	resultBox.Disable()

	statusLabel := widget.NewLabel("")

	generateBtn := widget.NewButtonWithIcon("Generate Ideas", theme.HelpIcon(), func() {
		if apiKey == "" {
			dialog.ShowInformation("API Key Required", "Please configure your Gemini API Key in the Settings tab or set GEMINI_API_KEY environment variable.", w)
			return
		}
		statusLabel.SetText("Generating sub-thoughts with Gemini...")
		go func() {
			out, err := ai.GenerateIdeasWithKey(promptEntry.Text, apiKey)
			if err != nil {
				statusLabel.SetText("Error: " + err.Error())
				return
			}
			statusLabel.SetText("Generation successful!")
			resultBox.SetText(out)
		}()
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle("Ask Gemini AI for Creative Sub-Thoughts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Prompt:"),
		promptEntry,
		container.NewHBox(generateBtn, statusLabel),
		widget.NewLabel("AI Suggestions:"),
		container.NewVScroll(resultBox),
	)

	d := dialog.NewCustomConfirm("AI Sub-Thought Generator", "Save as Sub-Thought", "Close", content, func(save bool) {
		if save && strings.TrimSpace(resultBox.Text) != "" {
			meta := &models.MetaData{
				ID:          models.NewID(),
				ParentID:    idea.ID,
				ParentType:  "Idea",
				Level:       1,
				Title:       "Gemini Suggestions",
				Description: strings.TrimSpace(resultBox.Text),
			}
			db.CreateMetaData(meta)
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(500, 480))
	d.Show()
}
