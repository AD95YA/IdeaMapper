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

// BuildMemoriesList constructs the timeline memories view with search, date grouping and cards
func BuildMemoriesList(a fyne.App, w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()
	scroll := container.NewVScroll(vbox)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search memories by title or details...")

	onlyImportantCheck := widget.NewCheck("★ Starred Only", nil)

	var refreshList func()

	refreshList = func() {
		vbox.Objects = nil
		memories, err := db.GetAllMemories()
		if err != nil {
			vbox.Add(widget.NewLabel("Failed to load memories: " + err.Error()))
			return
		}

		query := strings.ToLower(strings.TrimSpace(searchEntry.Text))
		onlyImp := onlyImportantCheck.Checked

		var currentYear int
		displayCount := 0

		for _, m := range memories {
			mem := m
			if onlyImp && !mem.IsImportant {
				continue
			}
			if query != "" {
				matches := strings.Contains(strings.ToLower(mem.Title), query) ||
					strings.Contains(strings.ToLower(mem.Description), query)
				if !matches {
					continue
				}
			}

			displayCount++

			year := mem.Date.Year()
			if year != currentYear {
				currentYear = year
				yearLabel := widget.NewLabelWithStyle(fmt.Sprintf("─── %d ───", year), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
				vbox.Add(yearLabel)
			}

			card := buildMemoryCard(a, w, mem, refreshList)
			vbox.Add(card)
		}

		if displayCount == 0 {
			vbox.Add(widget.NewLabelWithStyle("No memories match your filter.\nClick '+ Create Memory' to record a moment.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		}

		scroll.Refresh()
	}

	searchEntry.OnChanged = func(_ string) {
		refreshList()
	}
	onlyImportantCheck.OnChanged = func(_ bool) {
		refreshList()
	}

	addBtn := widget.NewButtonWithIcon("Create Memory", theme.ContentAddIcon(), func() {
		ShowCreateMemoryModal(w, refreshList)
	})

	topBar := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("My Memories & Timeline", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layoutSpacer(),
			addBtn,
		),
		container.NewBorder(nil, nil, nil, onlyImportantCheck, searchEntry),
	)

	refreshList()

	return container.NewBorder(topBar, nil, nil, nil, scroll)
}

func buildMemoryCard(a fyne.App, w fyne.Window, m *models.Memory, onRefresh func()) fyne.CanvasObject {
	starPrefix := "☆ "
	if m.IsImportant {
		starPrefix = "★ "
	}

	dateStr := m.Date.Format("Jan 02")
	titleLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s%s  •  %s", starPrefix, dateStr, m.Title), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	descText := m.Description
	if len(descText) > 120 {
		descText = descText[:117] + "..."
	}
	descLabel := widget.NewLabel(descText)

	metaList, _ := db.GetMetaDataForParent(m.ID, "Memory")
	links, _ := db.GetInterLinksForSource(m.ID)
	infoText := fmt.Sprintf("%d notes  •  %d connections", len(metaList), len(links))
	infoLabel := widget.NewLabelWithStyle(infoText, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	viewBtn := widget.NewButtonWithIcon("View", theme.FolderOpenIcon(), func() {
		ShowMemoryDetail(a, w, m.ID, func() {
			w.SetContent(BuildMainTabs(a, w))
		})
	})

	starToggleBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		m.IsImportant = !m.IsImportant
		db.UpdateMemory(m)
		onRefresh()
	})
	if !m.IsImportant {
		starToggleBtn.SetIcon(theme.CancelIcon())
	}

	delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Delete Memory", "Delete '"+m.Title+"'?", func(ok bool) {
			if ok {
				db.DeleteMemory(m.ID)
				onRefresh()
			}
		}, w)
	})

	actionsRow := container.NewHBox(viewBtn, starToggleBtn, delBtn)
	bottomRow := container.NewBorder(nil, nil, infoLabel, actionsRow)

	cardContent := container.NewVBox(titleLabel)
	if descText != "" {
		cardContent.Add(descLabel)
	}
	cardContent.Add(bottomRow)
	cardContent.Add(widget.NewSeparator())

	return cardContent
}

// ShowCreateMemoryModal prompts to create a new Memory
func ShowCreateMemoryModal(w fyne.Window, onSave func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Memory Title")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Describe what happened...")

	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))
	dateEntry.SetPlaceHolder("YYYY-MM-DD")

	impCheck := widget.NewCheck("Mark as Important (★)", nil)
	impCheck.SetChecked(true)

	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
		widget.NewFormItem("Date (YYYY-MM-DD)", dateEntry),
		widget.NewFormItem("Importance", impCheck),
	)

	d := dialog.NewCustomConfirm("New Memory", "Save", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateEntry.Text))
			if err != nil {
				parsedDate = time.Now()
			}

			mem := &models.Memory{
				ID:          models.NewID(),
				Title:       strings.TrimSpace(titleEntry.Text),
				Description: strings.TrimSpace(descEntry.Text),
				Date:        parsedDate,
				IsImportant: impCheck.Checked,
			}
			db.CreateMemory(mem)
			onSave()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 320))
	d.Show()
}
