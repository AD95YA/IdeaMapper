package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// BuildMainTabs creates the primary 4 bottom tabs of IdeaMapper with swipe navigation and top header
func BuildMainTabs(a fyne.App, w fyne.Window) fyne.CanvasObject {
	nav := GetNavManager(w)

	ideasScreen := BuildIdeasList(a, w)
	memoriesScreen := BuildMemoriesList(a, w)
	linksScreen := BuildLinksView(a, w)
	settingsScreen := BuildSettings(a, w)

	tabNames := []string{"💡 Ideas", "💭 Memories", "🧠 Mind Map", "⚙ Settings"}
	tabs := container.NewAppTabs(
		container.NewTabItem(tabNames[0], ideasScreen),
		container.NewTabItem(tabNames[1], memoriesScreen),
		container.NewTabItem(tabNames[2], linksScreen),
		container.NewTabItem(tabNames[3], settingsScreen),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	tabs.OnSelected = func(item *container.TabItem) {
		nav.SetActiveRootCrumb(item.Text)
	}

	swipeable := NewSwipeableTabs(tabs, func(newIndex int) {
		if newIndex >= 0 && newIndex < len(tabNames) {
			nav.SetActiveRootCrumb(tabNames[newIndex])
		}
	})

	nav.SetOnTabSelected(func(index int) {
		if index >= 0 && index < len(tabs.Items) {
			tabs.SelectIndex(index)
		}
	})

	nav.SetRootView(swipeable)
	return nav.BuildMasterLayout()
}
