package ui

import (
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NavEntry represents an entry in the navigation stack
type NavEntry struct {
	Title       string
	Breadcrumbs []string
	View        fyne.CanvasObject
	OnBack      func()
}

// NavigationManager coordinates hierarchical view transitions, top bar breadcrumbs, and back history
type NavigationManager struct {
	mu                  sync.Mutex
	window              fyne.Window
	stack               []NavEntry
	rootView            fyne.CanvasObject
	contentContainer    *fyne.Container
	headerContainer     *fyne.Container
	backBtn             *widget.Button
	homeBtn             *widget.Button
	breadcrumbBox       *fyne.Container
	breadcrumbScroll    *container.Scroll
	activeRootCrumb     string
	onTabSelected       func(int)
}

var (
	globalNavInstance *NavigationManager
	navOnce           sync.Once
)

// GetNavManager returns the active singleton NavigationManager
func GetNavManager(w fyne.Window) *NavigationManager {
	navOnce.Do(func() {
		globalNavInstance = NewNavigationManager(w)
	})
	if w != nil && globalNavInstance.window == nil {
		globalNavInstance.window = w
		globalNavInstance.setupKeyInterception()
	}
	return globalNavInstance
}

// ResetNavManagerForTest allows tests to instantiate a fresh NavigationManager
func ResetNavManagerForTest(w fyne.Window) *NavigationManager {
	globalNavInstance = NewNavigationManager(w)
	return globalNavInstance
}

// NewNavigationManager initializes a NavigationManager
func NewNavigationManager(w fyne.Window) *NavigationManager {
	nm := &NavigationManager{
		window:          w,
		activeRootCrumb: "💡 Ideas",
	}

	nm.breadcrumbBox = container.NewHBox()
	nm.breadcrumbScroll = container.NewHScroll(nm.breadcrumbBox)

	nm.backBtn = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		nm.Pop()
	})
	nm.backBtn.Importance = widget.LowImportance

	nm.homeBtn = widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		nm.ResetToRoot()
	})
	nm.homeBtn.Importance = widget.LowImportance

	nm.headerContainer = container.NewBorder(
		nil, nil,
		container.NewHBox(nm.backBtn, nm.homeBtn, widget.NewSeparator()),
		nil,
		nm.breadcrumbScroll,
	)

	nm.contentContainer = container.NewStack()
	nm.updateHeader()

	if w != nil {
		nm.setupKeyInterception()
	}

	return nm
}

func (nm *NavigationManager) setupKeyInterception() {
	if nm.window == nil {
		return
	}
	// Intercept Android hardware back / desktop Escape
	nm.window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			if nm.CanPop() {
				nm.Pop()
			}
		}
	})
}

// SetRootView sets the base 4-tab container
func (nm *NavigationManager) SetRootView(root fyne.CanvasObject) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.rootView = root
	if len(nm.stack) == 0 {
		nm.contentContainer.Objects = []fyne.CanvasObject{root}
		nm.contentContainer.Refresh()
	}
	nm.updateHeader()
}

// SetActiveRootCrumb updates breadcrumbs when user changes root tabs
func (nm *NavigationManager) SetActiveRootCrumb(title string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.activeRootCrumb = title
	if len(nm.stack) == 0 {
		nm.updateHeader()
	}
}

// SetOnTabSelected registers a callback for switching root tabs
func (nm *NavigationManager) SetOnTabSelected(fn func(int)) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.onTabSelected = fn
}

// BuildMasterLayout wraps header and content container into a unified canvas
func (nm *NavigationManager) BuildMasterLayout() fyne.CanvasObject {
	return container.NewBorder(nm.headerContainer, nil, nil, nil, nm.contentContainer)
}

// Push navigates forward to a new child/detail view
func (nm *NavigationManager) Push(entry NavEntry) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	nm.stack = append(nm.stack, entry)
	nm.contentContainer.Objects = []fyne.CanvasObject{entry.View}
	nm.contentContainer.Refresh()
	nm.updateHeader()
}

// Pop steps back one level in the hierarchy
func (nm *NavigationManager) Pop() bool {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if len(nm.stack) == 0 {
		return false
	}

	lastIdx := len(nm.stack) - 1
	popped := nm.stack[lastIdx]
	nm.stack = nm.stack[:lastIdx]

	if popped.OnBack != nil {
		popped.OnBack()
	}

	if len(nm.stack) > 0 {
		current := nm.stack[len(nm.stack)-1]
		nm.contentContainer.Objects = []fyne.CanvasObject{current.View}
	} else if nm.rootView != nil {
		nm.contentContainer.Objects = []fyne.CanvasObject{nm.rootView}
	}

	nm.contentContainer.Refresh()
	nm.updateHeader()
	return true
}

// PopToLevel pops the stack back to a specific breadcrumb index
func (nm *NavigationManager) PopToLevel(depth int) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if depth < 0 {
		depth = 0
	}
	if depth >= len(nm.stack) {
		return
	}

	// Trigger onBack for discarded entries
	for i := len(nm.stack) - 1; i >= depth; i-- {
		if nm.stack[i].OnBack != nil {
			nm.stack[i].OnBack()
		}
	}

	nm.stack = nm.stack[:depth]

	if len(nm.stack) > 0 {
		current := nm.stack[len(nm.stack)-1]
		nm.contentContainer.Objects = []fyne.CanvasObject{current.View}
	} else if nm.rootView != nil {
		nm.contentContainer.Objects = []fyne.CanvasObject{nm.rootView}
	}

	nm.contentContainer.Refresh()
	nm.updateHeader()
}

// UpdateCurrentView updates the view of the current top stack entry in-place
func (nm *NavigationManager) UpdateCurrentView(view fyne.CanvasObject) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if len(nm.stack) > 0 {
		nm.stack[len(nm.stack)-1].View = view
		nm.contentContainer.Objects = []fyne.CanvasObject{view}
		nm.contentContainer.Refresh()
	}
}

// CurrentEntry returns the active top NavEntry if stack is not empty
func (nm *NavigationManager) CurrentEntry() (NavEntry, bool) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	if len(nm.stack) == 0 {
		return NavEntry{}, false
	}
	return nm.stack[len(nm.stack)-1], true
}

// ResetToRoot jumps all the way back to the root tab dashboard
func (nm *NavigationManager) ResetToRoot() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for i := len(nm.stack) - 1; i >= 0; i-- {
		if nm.stack[i].OnBack != nil {
			nm.stack[i].OnBack()
		}
	}
	nm.stack = nil

	if nm.rootView != nil {
		nm.contentContainer.Objects = []fyne.CanvasObject{nm.rootView}
		nm.contentContainer.Refresh()
	}
	nm.updateHeader()
}

// CanPop returns true if there are sub-levels to navigate back from
func (nm *NavigationManager) CanPop() bool {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	return len(nm.stack) > 0
}

// Depth returns the current stack depth
func (nm *NavigationManager) Depth() int {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	return len(nm.stack)
}

func (nm *NavigationManager) updateHeader() {
	nm.breadcrumbBox.Objects = nil

	if len(nm.stack) == 0 {
		nm.backBtn.Disable()
		rootLabel := widget.NewLabelWithStyle("Mendu: "+nm.activeRootCrumb, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		nm.breadcrumbBox.Add(rootLabel)
	} else {
		nm.backBtn.Enable()
		nm.homeBtn.Enable()

		// Add Home root crumb
		homeCrumb := widget.NewButtonWithIcon(nm.activeRootCrumb, theme.HomeIcon(), func() {
			nm.ResetToRoot()
		})
		homeCrumb.Importance = widget.LowImportance
		nm.breadcrumbBox.Add(homeCrumb)

		currentEntry := nm.stack[len(nm.stack)-1]
		crumbs := currentEntry.Breadcrumbs
		if len(crumbs) == 0 && currentEntry.Title != "" {
			crumbs = []string{currentEntry.Title}
		}

		for idx, crumb := range crumbs {
			crumbIdx := idx
			nm.breadcrumbBox.Add(widget.NewLabel("›"))
			truncated := crumb
			if len(truncated) > 20 {
				truncated = strings.TrimSpace(truncated[:18]) + "…"
			}

			// If last crumb, render bold label; if ancestor crumb, render clickable jump button
			if idx == len(crumbs)-1 {
				leafLabel := widget.NewLabelWithStyle(truncated, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				nm.breadcrumbBox.Add(leafLabel)
			} else {
				crumbBtn := widget.NewButton(truncated, func() {
					// Pop to this level
					targetDepth := crumbIdx + 1
					if targetDepth < len(nm.stack) {
						nm.PopToLevel(targetDepth)
					}
				})
				crumbBtn.Importance = widget.LowImportance
				nm.breadcrumbBox.Add(crumbBtn)
			}
		}
	}

	nm.breadcrumbBox.Refresh()
	nm.breadcrumbScroll.Refresh()
}
