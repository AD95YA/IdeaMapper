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

// BuildIdeasList constructs the main Ideas dashboard view with search, filtering and cards
func BuildIdeasList(a fyne.App, w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()
	scroll := container.NewVScroll(vbox)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search ideas by title or content...")

	onlyImportantCheck := widget.NewCheck("★ Starred Only", nil)

	var refreshList func()

	refreshList = func() {
		vbox.Objects = nil
		ideas, err := db.GetAllIdeas()
		if err != nil {
			vbox.Add(widget.NewLabel("Failed to load ideas: " + err.Error()))
			return
		}

		query := strings.ToLower(strings.TrimSpace(searchEntry.Text))
		onlyImp := onlyImportantCheck.Checked

		displayCount := 0
		for _, idea := range ideas {
			if onlyImp && !idea.IsImportant {
				continue
			}
			if query != "" {
				matches := strings.Contains(strings.ToLower(idea.Title), query) ||
					strings.Contains(strings.ToLower(idea.Description), query)
				if !matches {
					continue
				}
			}

			displayCount++
			card := buildIdeaCard(a, w, idea, refreshList)
			vbox.Add(card)
		}

		if displayCount == 0 {
			vbox.Add(widget.NewLabelWithStyle("No ideas match your filter.\nClick '+ Create Idea' to record a new thought.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		}

		scroll.Refresh()
	}

	searchEntry.OnChanged = func(_ string) {
		refreshList()
	}
	onlyImportantCheck.OnChanged = func(_ bool) {
		refreshList()
	}

	addBtn := widget.NewButtonWithIcon("Create Idea", theme.ContentAddIcon(), func() {
		ShowCreateIdeaModal(w, refreshList)
	})

	topBar := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("My Ideas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layoutSpacer(),
			addBtn,
		),
		container.NewBorder(nil, nil, nil, onlyImportantCheck, searchEntry),
	)

	refreshList()

	return container.NewBorder(topBar, nil, nil, nil, scroll)
}

func buildIdeaCard(a fyne.App, w fyne.Window, idea *models.Idea, onRefresh func()) fyne.CanvasObject {
	starPrefix := "☆ "
	if idea.IsImportant {
		starPrefix = "★ "
	}

	titleLabel := widget.NewLabelWithStyle(starPrefix+idea.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	descText := idea.Description
	if len(descText) > 120 {
		descText = descText[:117] + "..."
	}
	descLabel := widget.NewLabel(descText)

	// Fetch count of sub-thoughts and links
	metaList, _ := db.GetMetaDataForParent(idea.ID, "Idea")
	links, _ := db.GetInterLinksForSource(idea.ID)
	infoText := fmt.Sprintf("%d sub-thoughts  •  %d connections", len(metaList), len(links))
	infoLabel := widget.NewLabelWithStyle(infoText, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	openBtn := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		ShowIdeaDetail(a, w, idea.ID, func() {
			w.SetContent(BuildMainTabs(a, w))
		})
	})

	starToggleBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		idea.IsImportant = !idea.IsImportant
		db.UpdateIdea(idea)
		onRefresh()
	})
	if !idea.IsImportant {
		starToggleBtn.SetIcon(theme.CancelIcon())
	}

	delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Delete Idea", "Delete '"+idea.Title+"'?", func(ok bool) {
			if ok {
				db.DeleteIdea(idea.ID)
				onRefresh()
			}
		}, w)
	})

	actionsRow := container.NewHBox(openBtn, starToggleBtn, delBtn)
	bottomRow := container.NewBorder(nil, nil, infoLabel, actionsRow)

	cardContent := container.NewVBox(titleLabel)
	if descText != "" {
		cardContent.Add(descLabel)
	}
	cardContent.Add(bottomRow)
	cardContent.Add(widget.NewSeparator())

	return cardContent
}

// ShowCreateIdeaModal prompts the user to create a new Idea
func ShowCreateIdeaModal(w fyne.Window, onSave func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Idea Title")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description / Thoughts")

	impCheck := widget.NewCheck("Mark as Important (★)", nil)

	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
		widget.NewFormItem("Importance", impCheck),
	)

	d := dialog.NewCustomConfirm("New Idea", "Save", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			idea := &models.Idea{
				ID:          models.NewID(),
				Title:       strings.TrimSpace(titleEntry.Text),
				Description: strings.TrimSpace(descEntry.Text),
				IsImportant: impCheck.Checked,
			}
			db.CreateIdea(idea)
			onSave()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 280))
	d.Show()
}
