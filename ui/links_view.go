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
)

// BuildLinksView creates the Mind Map explorer tab with Box-Drawing Tree, Radial Orbit, and Connection List views
func BuildLinksView(a fyne.App, w fyne.Window) fyne.CanvasObject {
	bodyStack := container.NewStack()

	modes := []string{"🌳 Box-Drawing Tree", "🪐 Radial Orbit", "🔗 Connections List"}
	var activeMode string = modes[0]

	var renderActiveMode func()

	renderActiveMode = func() {
		bodyStack.Objects = nil
		switch activeMode {
		case "🌳 Box-Drawing Tree":
			bodyStack.Objects = []fyne.CanvasObject{
				BuildAsciiMindMapView(a, w, func(nodeType, id, title string) {
					if nodeType == "Idea" {
						ShowIdeaDetail(a, w, id, renderActiveMode)
					} else if nodeType == "Memory" {
						ShowMemoryDetail(a, w, id, renderActiveMode)
					}
				}),
			}
		case "🪐 Radial Orbit":
			bodyStack.Objects = []fyne.CanvasObject{BuildRadialOrbitView(a, w)}
		case "🔗 Connections List":
			bodyStack.Objects = []fyne.CanvasObject{buildClassicLinksList(a, w)}
		}
		bodyStack.Refresh()
	}

	modeSelect := widget.NewSelect(modes, func(s string) {
		activeMode = s
		renderActiveMode()
	})
	modeSelect.SetSelected(activeMode)

	topBar := container.NewBorder(
		nil, nil,
		widget.NewLabelWithStyle("Mind Map Mode:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		modeSelect,
	)

	renderActiveMode()

	return container.NewBorder(topBar, nil, nil, nil, bodyStack)
}

func buildClassicLinksList(a fyne.App, w fyne.Window) fyne.CanvasObject {
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

			btnRow := container.NewHBox(openSourceBtn, openTargetBtn, delBtn)
			card := container.NewVBox(summaryLabel, btnRow, widget.NewSeparator())
			vbox.Add(card)
		}

		if count == 0 {
			vbox.Add(widget.NewLabelWithStyle("No connections found.\nUse 'Link to...' in any detail view to connect ideas & memories.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		}
		scroll.Refresh()
	}

	searchEntry.OnChanged = func(_ string) {
		refreshList()
	}

	refreshList()

	return container.NewBorder(searchEntry, nil, nil, nil, scroll)
}
