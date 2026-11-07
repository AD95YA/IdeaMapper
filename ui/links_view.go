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

// BuildLinksView creates the Mind Map & InterLinks explorer tab
func BuildLinksView(a fyne.App, w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()
	scroll := container.NewVScroll(vbox)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search connections by title or relationship...")

	var refreshList func()

	refreshList = func() {
		vbox.Objects = nil
		links, err := db.GetAllInterLinks()
		if err != nil {
			vbox.Add(widget.NewLabel("Failed to load connections: " + err.Error()))
			return
		}

		query := strings.ToLower(strings.TrimSpace(searchEntry.Text))

		count := 0
		for _, l := range links {
			link := l

			var sourceTitle, targetTitle string
			if link.SourceType == "Idea" {
				if idea, _ := db.GetIdeaByID(link.SourceID); idea != nil {
					sourceTitle = idea.Title
				}
			} else {
				if mem, _ := db.GetMemoryByID(link.SourceID); mem != nil {
					sourceTitle = mem.Title
				}
			}

			if link.TargetType == "Idea" {
				if idea, _ := db.GetIdeaByID(link.TargetID); idea != nil {
					targetTitle = idea.Title
				}
			} else {
				if mem, _ := db.GetMemoryByID(link.TargetID); mem != nil {
					targetTitle = mem.Title
				}
			}

			if sourceTitle == "" {
				sourceTitle = "(Unknown)"
			}
			if targetTitle == "" {
				targetTitle = "(Unknown)"
			}

			// Filter check
			if query != "" {
				matches := strings.Contains(strings.ToLower(sourceTitle), query) ||
					strings.Contains(strings.ToLower(targetTitle), query) ||
					strings.Contains(strings.ToLower(link.Name), query)
				if !matches {
					continue
				}
			}

			count++

			relName := link.Name
			if relName == "" {
				relName = "Connected to"
			}

			summaryLabel := widget.NewLabelWithStyle(
				fmt.Sprintf("• [%s] %s  ───(%s)───➔  [%s] %s",
					link.SourceType, sourceTitle, relName, link.TargetType, targetTitle),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			openSourceBtn := widget.NewButton("Open "+link.SourceType, func() {
				if link.SourceType == "Idea" {
					ShowIdeaDetail(a, w, link.SourceID, func() {
						w.SetContent(BuildMainTabs(a, w))
					})
				} else {
					ShowMemoryDetail(a, w, link.SourceID, func() {
						w.SetContent(BuildMainTabs(a, w))
					})
				}
			})

			openTargetBtn := widget.NewButton("Open "+link.TargetType, func() {
				if link.TargetType == "Idea" {
					ShowIdeaDetail(a, w, link.TargetID, func() {
						w.SetContent(BuildMainTabs(a, w))
					})
				} else {
					ShowMemoryDetail(a, w, link.TargetID, func() {
						w.SetContent(BuildMainTabs(a, w))
					})
				}
			})

			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Delete Link", "Disconnect this link?", func(ok bool) {
					if ok {
						db.DeleteInterLink(link.ID)
						refreshList()
					}
				}, w)
			})

			actionsRow := container.NewHBox(openSourceBtn, openTargetBtn, delBtn)
			cardContent := container.NewVBox(summaryLabel)

			if link.LinkMetaData != nil && link.LinkMetaData.Title != "" {
				noteLabel := widget.NewLabelWithStyle("   Link Note: "+link.LinkMetaData.Title, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
				cardContent.Add(noteLabel)
			}

			cardContent.Add(actionsRow)
			cardContent.Add(widget.NewSeparator())

			vbox.Add(cardContent)
		}

		if count == 0 {
			vbox.Add(widget.NewLabelWithStyle("No matching connections found.\nUse '+ New Connection' or open an Idea/Memory to link ideas and thoughts together.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		}

		scroll.Refresh()
	}

	searchEntry.OnChanged = func(_ string) {
		refreshList()
	}

	addLinkBtn := widget.NewButtonWithIcon("New Connection", theme.ContentAddIcon(), func() {
		// Can pick any source and target
		ShowGlobalCreateLinkModal(w, refreshList)
	})

	topBar := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Mind Map & Connections", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layoutSpacer(),
			addLinkBtn,
		),
		searchEntry,
	)

	refreshList()

	return container.NewBorder(topBar, nil, nil, nil, scroll)
}

// ShowGlobalCreateLinkModal allows creating a connection picking both Source and Target
func ShowGlobalCreateLinkModal(w fyne.Window, onSaved func()) {
	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()

	type ItemChoice struct {
		ID    string
		Type  string
		Label string
	}
	var choices []ItemChoice
	var names []string

	for _, i := range ideas {
		lbl := "[Idea] " + i.Title
		choices = append(choices, ItemChoice{ID: i.ID, Type: "Idea", Label: lbl})
		names = append(names, lbl)
	}
	for _, m := range memories {
		lbl := "[Memory] " + m.Title
		choices = append(choices, ItemChoice{ID: m.ID, Type: "Memory", Label: lbl})
		names = append(names, lbl)
	}

	if len(choices) < 2 {
		dialog.ShowInformation("Not Enough Items", "You need at least 2 ideas or memories to create a connection.", w)
		return
	}

	sourceSelect := widget.NewSelect(names, nil)
	sourceSelect.SetSelected(names[0])

	targetSelect := widget.NewSelect(names, nil)
	targetSelect.SetSelected(names[1])

	relNameEntry := widget.NewEntry()
	relNameEntry.SetPlaceHolder("Relationship (e.g. Inspires, Related to, Causes)")
	relNameEntry.SetText("Related to")

	noteEntry := widget.NewEntry()
	noteEntry.SetPlaceHolder("Optional Connection Note")

	form := widget.NewForm(
		widget.NewFormItem("From (Source)", sourceSelect),
		widget.NewFormItem("To (Target)", targetSelect),
		widget.NewFormItem("Relationship", relNameEntry),
		widget.NewFormItem("Link Note", noteEntry),
	)

	d := dialog.NewCustomConfirm("Create Mind Map Connection", "Create Connection", "Cancel", form, func(save bool) {
		if save {
			if sourceSelect.Selected == targetSelect.Selected {
				dialog.ShowError(fmt.Errorf("source and target cannot be the same item"), w)
				return
			}

			var src, tgt *ItemChoice
			for i := range choices {
				if choices[i].Label == sourceSelect.Selected {
					src = &choices[i]
				}
				if choices[i].Label == targetSelect.Selected {
					tgt = &choices[i]
				}
			}

			if src == nil || tgt == nil {
				return
			}

			link := &models.InterLink{
				ID:         models.NewID(),
				SourceID:   src.ID,
				SourceType: src.Type,
				TargetID:   tgt.ID,
				TargetType: tgt.Type,
				Name:       strings.TrimSpace(relNameEntry.Text),
			}

			if err := db.CreateInterLink(link); err != nil {
				dialog.ShowError(err, w)
				return
			}

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

	d.Resize(fyne.NewSize(420, 320))
	d.Show()
}
