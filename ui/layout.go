package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// BuildMainTabs creates the primary 4 bottom tabs of IdeaMapper
func BuildMainTabs(a fyne.App, w fyne.Window) fyne.CanvasObject {
	ideasScreen := BuildIdeasList(a, w)
	memoriesScreen := BuildMemoriesList(a, w)
	linksScreen := BuildLinksView(a, w)
	settingsScreen := BuildSettings(a, w)

	tabs := container.NewAppTabs(
		container.NewTabItem("💡 Ideas", ideasScreen),
		container.NewTabItem("💭 Memories", memoriesScreen),
		container.NewTabItem("🧠 Mind Map", linksScreen),
		container.NewTabItem("⚙ Settings", settingsScreen),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	return tabs
}
