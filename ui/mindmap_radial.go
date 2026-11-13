package ui

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
)

// LevelAccentColor returns depth-graded accent colors for levels 1 through 6
func LevelAccentColor(level int) color.Color {
	switch level {
	case 1:
		return color.RGBA{R: 6, G: 182, B: 212, A: 255} // Cyan
	case 2:
		return color.RGBA{R: 245, G: 158, B: 11, A: 255} // Amber
	case 3:
		return color.RGBA{R: 16, G: 185, B: 129, A: 255} // Emerald
	case 4:
		return color.RGBA{R: 139, G: 92, B: 246, A: 255} // Violet
	case 5:
		return color.RGBA{R: 244, G: 63, B: 94, A: 255} // Rose
	default:
		return color.RGBA{R: 148, G: 163, B: 184, A: 255} // Slate
	}
}

// BuildRadialOrbitView creates the active-orbit radial mind map visualization
func BuildRadialOrbitView(a fyne.App, w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()

	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()

	if len(ideas) == 0 && len(memories) == 0 {
		return widget.NewLabelWithStyle("No topics available for radial orbit rendering.\nCreate an Idea or Memory to view orbits.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	}

	// Topic selector options
	var options []string
	idMap := make(map[string]struct {
		id   string
		typ  string
		name string
	})

	for _, id := range ideas {
		opt := "💡 " + id.Title
		options = append(options, opt)
		idMap[opt] = struct {
			id   string
			typ  string
			name string
		}{id: id.ID, typ: "Idea", name: id.Title}
	}
	for _, m := range memories {
		opt := "💭 " + m.Title
		options = append(options, opt)
		idMap[opt] = struct {
			id   string
			typ  string
			name string
		}{id: m.ID, typ: "Memory", name: m.Title}
	}

	selectedOpt := options[0]
	orbitContainer := container.NewVBox()

	renderOrbit := func(opt string) {
		orbitContainer.Objects = nil
		target := idMap[opt]

		// Fetch direct children & connections
		metaList, _ := db.GetMetaDataForParent(target.id, target.typ)
		links, _ := db.GetInterLinksForSource(target.id)

		// Center card
		centerHeader := widget.NewLabelWithStyle(fmt.Sprintf("🪐 Orbit Center: %s (%s)", target.name, target.typ), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		openDetailBtn := widget.NewButtonWithIcon("Open Detail View", theme.FolderOpenIcon(), func() {
			if target.typ == "Idea" {
				ShowIdeaDetail(a, w, target.id, nil)
			} else {
				ShowMemoryDetail(a, w, target.id, nil)
			}
		})

		orbitContainer.Add(container.NewHBox(layoutSpacer(), centerHeader, openDetailBtn, layoutSpacer()))
		orbitContainer.Add(widget.NewSeparator())

		// Sub-Thoughts (Level 1..6)
		if len(metaList) > 0 {
			orbitContainer.Add(widget.NewLabelWithStyle("Satellite Sub-Thoughts:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			for idx, m := range metaList {
				sub := m
				angleDeg := float64(idx) * (360.0 / float64(len(metaList)))
				badge := fmt.Sprintf("[L%d • %0.0f°] ", sub.Level, angleDeg)

				subBtn := widget.NewButtonWithIcon(badge+sub.Title, theme.NavigateNextIcon(), func() {
					drillView := BuildFocusedOrbitView(w, sub, func() {
						w.SetContent(BuildMainTabs(a, w))
					})
					GetNavManager(w).Push(NavEntry{
						Title:       sub.Title,
						Breadcrumbs: []string{target.name, fmt.Sprintf("L%d: %s", sub.Level, sub.Title)},
						View:        drillView,
					})
				})
				subBtn.Importance = widget.MediumImportance

				metaInfo := fmt.Sprintf("%d nested sub-items", len(sub.SubMetaData))
				if len(sub.PhotoPaths) > 0 {
					metaInfo += fmt.Sprintf(" • 📷 %d attachments", len(sub.PhotoPaths))
				}
				row := container.NewBorder(nil, nil, nil, widget.NewLabel(metaInfo), subBtn)
				orbitContainer.Add(row)
			}
			orbitContainer.Add(widget.NewSeparator())
		}

		// InterLink Connections
		if len(links) > 0 {
			orbitContainer.Add(widget.NewLabelWithStyle("Radiating InterLinks:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			for _, l := range links {
				link := l
				relName := link.Name
				if relName == "" {
					relName = "Connected to"
				}

				var targetName string
				if link.TargetType == "Idea" {
					if targetIdea, _ := db.GetIdeaByID(link.TargetID); targetIdea != nil {
						targetName = targetIdea.Title
					}
				} else {
					if targetMem, _ := db.GetMemoryByID(link.TargetID); targetMem != nil {
						targetName = targetMem.Title
					}
				}
				if targetName == "" {
					targetName = "(Unknown)"
				}

				linkBtn := widget.NewButton(fmt.Sprintf("──(%s)──➔ [%s] %s", relName, link.TargetType, targetName), func() {
					if link.TargetType == "Idea" {
						ShowIdeaDetail(a, w, link.TargetID, nil)
					} else {
						ShowMemoryDetail(a, w, link.TargetID, nil)
					}
				})
				orbitContainer.Add(linkBtn)
			}
		}

		if len(metaList) == 0 && len(links) == 0 {
			orbitContainer.Add(widget.NewLabelWithStyle("No satellite sub-thoughts or radiating links yet.\nAdd sub-thoughts in Detail View to expand this orbit.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		}
		orbitContainer.Refresh()
	}

	selectWidget := widget.NewSelect(options, func(opt string) {
		selectedOpt = opt
		renderOrbit(opt)
	})
	selectWidget.SetSelected(selectedOpt)

	topBar := container.NewBorder(
		nil, nil,
		widget.NewLabelWithStyle("Select Focus Topic:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		selectWidget,
	)

	vbox.Add(topBar)
	vbox.Add(orbitContainer)

	return container.NewVScroll(vbox)
}

// CalculateRadialPositions computes (x,y) coordinates on an ellipse to avoid text overlap
func CalculateRadialPositions(centerX, centerY, radiusX, radiusY float64, count int) []fyne.Position {
	if count <= 0 {
		return nil
	}
	positions := make([]fyne.Position, count)
	for i := 0; i < count; i++ {
		angle := (float64(i) * 2 * math.Pi) / float64(count)
		x := centerX + radiusX*math.Cos(angle)
		y := centerY + radiusY*math.Sin(angle)
		positions[i] = fyne.NewPos(float32(x), float32(y))
	}
	return positions
}
