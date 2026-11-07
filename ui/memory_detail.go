package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

// ShowMemoryDetail renders the full detail view for a Memory
func ShowMemoryDetail(a fyne.App, w fyne.Window, memoryID string, onBack func()) {
	var renderView func()

	renderView = func() {
		mem, err := db.GetMemoryByID(memoryID)
		if err != nil || mem == nil {
			dialog.ShowError(fmt.Errorf("memory not found or failed to load"), w)
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
		if mem.IsImportant {
			starIcon = theme.ConfirmIcon()
			starText = "★ Important"
		} else {
			starIcon = theme.CancelIcon()
			starText = "☆ Normal"
		}

		importantBtn := widget.NewButtonWithIcon(starText, starIcon, func() {
			mem.IsImportant = !mem.IsImportant
			if err := db.UpdateMemory(mem); err != nil {
				dialog.ShowError(err, w)
				return
			}
			renderView()
		})

		editBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
			showEditMemoryModal(w, mem, renderView)
		})

		deleteBtn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Delete Memory", "Are you sure you want to delete '"+mem.Title+"'?", func(ok bool) {
				if ok {
					db.DeleteMemory(mem.ID)
					if onBack != nil {
						onBack()
					}
				}
			}, w)
		})

		actionsBox := container.NewHBox(importantBtn, editBtn, deleteBtn)
		topBar := container.NewBorder(nil, nil, backBtn, actionsBox, widget.NewLabelWithStyle(mem.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

		// Body Content
		body := container.NewVBox()

		// Date Banner
		dateStr := mem.Date.Format("Monday, January 02, 2006")
		dateBanner := widget.NewLabelWithStyle("📅 Memory Date: "+dateStr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
		body.Add(dateBanner)

		// Description card
		descText := mem.Description
		if descText == "" {
			descText = "(No description provided)"
		}
		descCard := widget.NewCard("Description", "", widget.NewLabel(descText))
		body.Add(descCard)
		body.Add(widget.NewSeparator())

		// InterLinks Section
		linksSection := BuildInterLinksSection(a, w, mem.ID, "Memory", renderView, func(targetType string, targetID string) {
			if targetType == "Idea" {
				ShowIdeaDetail(a, w, targetID, renderView)
			} else {
				ShowMemoryDetail(a, w, targetID, renderView)
			}
		})
		body.Add(linksSection)
		body.Add(widget.NewSeparator())

		// MetaData Hierarchical Tree
		metaDataSection := BuildMetaDataTree(w, mem.ID, "Memory", renderView)
		body.Add(metaDataSection)

		scroll := container.NewVScroll(body)
		mainLayout := container.NewBorder(topBar, nil, nil, nil, scroll)

		w.SetContent(mainLayout)
	}

	renderView()
}

func showEditMemoryModal(w fyne.Window, mem *models.Memory, onSaved func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetText(mem.Title)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(mem.Description)

	dateEntry := widget.NewEntry()
	dateEntry.SetText(mem.Date.Format("2006-01-02"))
	dateEntry.SetPlaceHolder("YYYY-MM-DD")

	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
		widget.NewFormItem("Date (YYYY-MM-DD)", dateEntry),
	)

	d := dialog.NewCustomConfirm("Edit Memory", "Save", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateEntry.Text))
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid date format: use YYYY-MM-DD"), w)
				return
			}

			mem.Title = strings.TrimSpace(titleEntry.Text)
			mem.Description = strings.TrimSpace(descEntry.Text)
			mem.Date = parsedDate

			if err := db.UpdateMemory(mem); err != nil {
				dialog.ShowError(err, w)
				return
			}
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 280))
	d.Show()
}
