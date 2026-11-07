package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

// BuildInterLinksSection creates the UI section displaying connected links for an Idea or Memory
func BuildInterLinksSection(a fyne.App, w fyne.Window, currentID string, currentType string, onRefresh func(), onNavigate func(targetType string, targetID string)) fyne.CanvasObject {
	links, err := db.GetInterLinksForSource(currentID)
	if err != nil {
		return widget.NewLabel("Failed to load connections: " + err.Error())
	}

	content := container.NewVBox()

	addLinkBtn := widget.NewButtonWithIcon("Connect / Link", theme.ContentAddIcon(), func() {
		ShowCreateLinkModal(w, currentID, currentType, onRefresh)
	})

	content.Add(container.NewHBox(
		widget.NewLabelWithStyle("Connected Links (InterLinks)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layoutSpacer(),
		addLinkBtn,
	))

	if len(links) == 0 {
		content.Add(widget.NewLabelWithStyle("No links connected to this item.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		return content
	}

	for _, l := range links {
		link := l
		var otherID, otherType string
		var direction string

		if link.SourceID == currentID {
			otherID = link.TargetID
			otherType = link.TargetType
			direction = "➔ " + otherType + ": "
		} else {
			otherID = link.SourceID
			otherType = link.SourceType
			direction = "🡰 from " + otherType + ": "
		}

		// Fetch other item title
		var otherTitle string
		if otherType == "Idea" {
			if idea, err := db.GetIdeaByID(otherID); err == nil && idea != nil {
				otherTitle = idea.Title
			} else {
				otherTitle = "(Unknown Idea)"
			}
		} else {
			if mem, err := db.GetMemoryByID(otherID); err == nil && mem != nil {
				otherTitle = mem.Title
			} else {
				otherTitle = "(Unknown Memory)"
			}
		}

		relationName := link.Name
		if relationName == "" {
			relationName = "Related to"
		}

		titleText := fmt.Sprintf("[%s] %s%s", relationName, direction, otherTitle)
		linkTitleBtn := widget.NewButton(titleText, func() {
			if onNavigate != nil {
				onNavigate(otherType, otherID)
			}
		})

		delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Remove Link", "Disconnect this link?", func(ok bool) {
				if ok {
					db.DeleteInterLink(link.ID)
					onRefresh()
				}
			}, w)
		})

		row := container.NewBorder(nil, nil, nil, delBtn, linkTitleBtn)
		content.Add(row)

		if link.LinkMetaData != nil && link.LinkMetaData.Title != "" {
			noteLabel := widget.NewLabelWithStyle("    Note: "+link.LinkMetaData.Title+" - "+link.LinkMetaData.Description, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
			content.Add(noteLabel)
		}
	}

	return content
}

// ShowCreateLinkModal allows the user to choose another Idea or Memory to link with
func ShowCreateLinkModal(w fyne.Window, sourceID string, sourceType string, onSaved func()) {
	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()

	targetTypeSelect := widget.NewSelect([]string{"Idea", "Memory"}, nil)
	targetItemSelect := widget.NewSelect([]string{}, nil)

	type OptionInfo struct {
		ID   string
		Name string
	}
	var currentOptions []OptionInfo

	populateTargets := func(chosenType string) {
		currentOptions = nil
		var names []string

		if chosenType == "Idea" {
			for _, i := range ideas {
				if !(sourceType == "Idea" && i.ID == sourceID) {
					currentOptions = append(currentOptions, OptionInfo{ID: i.ID, Name: i.Title})
					names = append(names, i.Title)
				}
			}
		} else {
			for _, m := range memories {
				if !(sourceType == "Memory" && m.ID == sourceID) {
					currentOptions = append(currentOptions, OptionInfo{ID: m.ID, Name: m.Title})
					names = append(names, m.Title)
				}
			}
		}
		targetItemSelect.Options = names
		if len(names) > 0 {
			targetItemSelect.SetSelected(names[0])
		} else {
			targetItemSelect.ClearSelected()
		}
	}

	targetTypeSelect.OnChanged = func(val string) {
		populateTargets(val)
	}

	targetTypeSelect.SetSelected("Idea")

	relNameEntry := widget.NewEntry()
	relNameEntry.SetPlaceHolder("e.g. Inspired by, Causes, Sub-project, Related to")
	relNameEntry.SetText("Related to")

	noteEntry := widget.NewEntry()
	noteEntry.SetPlaceHolder("Optional Link Note / MetaData")

	form := widget.NewForm(
		widget.NewFormItem("Connect To", targetTypeSelect),
		widget.NewFormItem("Select Item", targetItemSelect),
		widget.NewFormItem("Relationship", relNameEntry),
		widget.NewFormItem("Link Note", noteEntry),
	)

	d := dialog.NewCustomConfirm("Create InterLink", "Link", "Cancel", form, func(save bool) {
		if save && targetItemSelect.Selected != "" {
			var chosenTargetID string
			for _, opt := range currentOptions {
				if opt.Name == targetItemSelect.Selected {
					chosenTargetID = opt.ID
					break
				}
			}

			if chosenTargetID == "" {
				dialog.ShowError(fmt.Errorf("no target item selected"), w)
				return
			}

			link := &models.InterLink{
				ID:         models.NewID(),
				SourceID:   sourceID,
				SourceType: sourceType,
				TargetID:   chosenTargetID,
				TargetType: targetTypeSelect.Selected,
				Name:       strings.TrimSpace(relNameEntry.Text),
			}

			if err := db.CreateInterLink(link); err != nil {
				dialog.ShowError(err, w)
				return
			}

			// If optional note is provided, create LinkMetaData
			if strings.TrimSpace(noteEntry.Text) != "" {
				meta := &models.MetaData{
					ID:         models.NewID(),
					ParentID:   link.ID,
					ParentType: "InterLink",
					Level:      1,
					Title:      strings.TrimSpace(noteEntry.Text),
				}
				db.CreateMetaData(meta)
			}

			onSaved()
		}
	}, w)

	d.Resize(fyne.NewSize(400, 320))
	d.Show()
}
