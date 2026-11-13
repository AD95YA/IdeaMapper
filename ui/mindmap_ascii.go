package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

// AsciiNode represents a node in the box-drawing tree
type AsciiNode struct {
	ID          string
	Type        string // "Idea", "Memory", "MetaData"
	Title       string
	Level       int
	IsImportant bool
	HasNotes    bool
	HasImages   bool
	IsIsolated  bool
	Children    []*AsciiNode
}

// BuildAsciiMindMapView builds the high-clarity Unicode box-drawing mind map tree view
func BuildAsciiMindMapView(a fyne.App, w fyne.Window, onFocus func(nodeType, id, title string)) fyne.CanvasObject {
	vbox := container.NewVBox()

	ideas, _ := db.GetAllIdeas()
	memories, _ := db.GetAllMemories()
	links, _ := db.GetAllInterLinks()

	// Track linked entity IDs
	linkedIDs := make(map[string]bool)
	for _, l := range links {
		linkedIDs[l.SourceID] = true
		linkedIDs[l.TargetID] = true
	}

	headerLabel := widget.NewLabelWithStyle("🌳 Mind Map: Box-Drawing Hierarchy", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	legendLabel := widget.NewLabelWithStyle("[★] Starred  •  [📝] Has Notes  •  [🖼️] Attachments  •  [⏳] Unlinked / Solo Thought", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	vbox.Add(headerLabel)
	vbox.Add(legendLabel)
	vbox.Add(widget.NewSeparator())

	if len(ideas) == 0 && len(memories) == 0 {
		vbox.Add(widget.NewLabelWithStyle("Your mind map is empty.\nCreate ideas or memories to populate the tree.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		return container.NewVScroll(vbox)
	}

	// 1. Ideas Section
	if len(ideas) > 0 {
		vbox.Add(widget.NewLabelWithStyle("💡 Ideas Root", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for i, idea := range ideas {
			isLast := (i == len(ideas)-1) && len(memories) == 0
			treeNode := buildIdeaAsciiNode(idea, linkedIDs)
			renderAsciiNode(vbox, a, w, treeNode, "", isLast, onFocus)
		}
		vbox.Add(widget.NewSeparator())
	}

	// 2. Memories Section
	if len(memories) > 0 {
		vbox.Add(widget.NewLabelWithStyle("💭 Memories Root", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for i, mem := range memories {
			isLast := (i == len(memories) - 1)
			treeNode := buildMemoryAsciiNode(mem, linkedIDs)
			renderAsciiNode(vbox, a, w, treeNode, "", isLast, onFocus)
		}
	}

	return container.NewVScroll(vbox)
}

func buildIdeaAsciiNode(idea *models.Idea, linkedIDs map[string]bool) *AsciiNode {
	metaList, _ := db.GetMetaDataForParent(idea.ID, "Idea")
	hasImg := false
	for _, m := range metaList {
		if len(m.PhotoPaths) > 0 {
			hasImg = true
			break
		}
	}

	isIsolated := !linkedIDs[idea.ID] && len(metaList) == 0

	node := &AsciiNode{
		ID:          idea.ID,
		Type:        "Idea",
		Title:       idea.Title,
		Level:       1,
		IsImportant: idea.IsImportant,
		HasNotes:    len(metaList) > 0,
		HasImages:   hasImg,
		IsIsolated:  isIsolated,
	}

	for _, m := range metaList {
		node.Children = append(node.Children, buildMetaDataAsciiNode(m))
	}
	return node
}

func buildMemoryAsciiNode(mem *models.Memory, linkedIDs map[string]bool) *AsciiNode {
	metaList, _ := db.GetMetaDataForParent(mem.ID, "Memory")
	hasImg := false
	for _, m := range metaList {
		if len(m.PhotoPaths) > 0 {
			hasImg = true
			break
		}
	}

	isIsolated := !linkedIDs[mem.ID] && len(metaList) == 0

	node := &AsciiNode{
		ID:          mem.ID,
		Type:        "Memory",
		Title:       mem.Title,
		Level:       1,
		IsImportant: mem.IsImportant,
		HasNotes:    len(metaList) > 0,
		HasImages:   hasImg,
		IsIsolated:  isIsolated,
	}

	for _, m := range metaList {
		node.Children = append(node.Children, buildMetaDataAsciiNode(m))
	}
	return node
}

func buildMetaDataAsciiNode(m *models.MetaData) *AsciiNode {
	node := &AsciiNode{
		ID:          m.ID,
		Type:        "MetaData",
		Title:       m.Title,
		Level:       m.Level,
		HasNotes:    m.Description != "",
		HasImages:   len(m.PhotoPaths) > 0,
		IsIsolated:  false,
	}
	for _, sub := range m.SubMetaData {
		node.Children = append(node.Children, buildMetaDataAsciiNode(sub))
	}
	return node
}

func renderAsciiNode(containerBox *fyne.Container, a fyne.App, w fyne.Window, node *AsciiNode, prefix string, isLast bool, onFocus func(string, string, string)) {
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	badges := ""
	if node.IsImportant {
		badges += "[★] "
	}
	if node.HasNotes {
		badges += "[📝] "
	}
	if node.HasImages {
		badges += "[🖼️] "
	}
	if node.IsIsolated {
		badges += "[⏳] "
	}

	lineText := fmt.Sprintf("%s%s%s%s", prefix, connector, badges, node.Title)
	lineLabel := widget.NewLabelWithStyle(lineText, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	focusBtn := widget.NewButtonWithIcon("Focus", theme.NavigateNextIcon(), func() {
		if onFocus != nil {
			onFocus(node.Type, node.ID, node.Title)
		} else {
			if node.Type == "Idea" {
				ShowIdeaDetail(a, w, node.ID, nil)
			} else if node.Type == "Memory" {
				ShowMemoryDetail(a, w, node.ID, nil)
			}
		}
	})
	focusBtn.Importance = widget.LowImportance

	row := container.NewHBox(lineLabel, focusBtn)
	containerBox.Add(row)

	var nextPrefix string
	if isLast {
		nextPrefix = prefix + "    "
	} else {
		nextPrefix = prefix + "│   "
	}

	for i, child := range node.Children {
		childLast := (i == len(node.Children)-1)
		renderAsciiNode(containerBox, a, w, child, nextPrefix, childLast, onFocus)
	}
}
