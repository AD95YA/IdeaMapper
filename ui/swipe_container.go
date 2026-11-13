package ui

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	// SwipeThresholdDP is the minimum horizontal drag distance in DP to trigger a tab change
	SwipeThresholdDP float32 = 40.0
	// EdgeSwipeZoneDP defines the edge margin for the Mind Map conflict guard
	EdgeSwipeZoneDP float32 = 24.0
	// HeaderSwipeZoneDP defines the top margin for header swipes
	HeaderSwipeZoneDP float32 = 55.0
	// MindMapTabIndex is the index of the Mind Map tab (0: Ideas, 1: Memories, 2: Mind Map, 3: Settings)
	MindMapTabIndex = 2
)

// EvaluateSwipeTransition calculates whether a drag gesture should trigger a tab switch
// Returns the target tab index and a boolean indicating if a transition occurred
func EvaluateSwipeTransition(startX, startY, width, totalDX, totalDY float32, currentTab int, totalTabs int) (int, bool) {
	absX := float32(math.Abs(float64(totalDX)))
	absY := float32(math.Abs(float64(totalDY)))

	// Must exceed horizontal threshold and be dominantly horizontal
	if absX < SwipeThresholdDP || absX < 1.5*absY {
		return currentTab, false
	}

	// Gesture Conflict Guard: On Mind Map tab, reserve central canvas for map panning
	if currentTab == MindMapTabIndex {
		isEdgeSwipe := startX < EdgeSwipeZoneDP || startX > (width-EdgeSwipeZoneDP)
		isHeaderSwipe := startY < HeaderSwipeZoneDP
		if !isEdgeSwipe && !isHeaderSwipe {
			// Center drag on Mind Map: suppress tab transition to preserve canvas panning
			return currentTab, false
		}
	}

	if totalDX < 0 {
		// Swiped Left -> Move forward (next tab)
		if currentTab < totalTabs-1 {
			return currentTab + 1, true
		}
	} else if totalDX > 0 {
		// Swiped Right -> Move backward (previous tab)
		if currentTab > 0 {
			return currentTab - 1, true
		}
	}

	return currentTab, false
}

// SwipeableTabs wraps an AppTabs container and adds horizontal gesture tab switching with conflict guards
type SwipeableTabs struct {
	widget.BaseWidget
	Tabs        *container.AppTabs
	startPos    fyne.Position
	totalDX     float32
	totalDY     float32
	isDragging  bool
	OnTabChange func(int)
}

// NewSwipeableTabs creates a gesture-enabled tab container
func NewSwipeableTabs(tabs *container.AppTabs, onTabChange func(int)) *SwipeableTabs {
	st := &SwipeableTabs{
		Tabs:        tabs,
		OnTabChange: onTabChange,
	}
	st.ExtendBaseWidget(st)
	return st
}

// CreateRenderer provides standard Fyne widget rendering
func (st *SwipeableTabs) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(st.Tabs)
}

// Dragged accumulates drag delta and tracks initial position
func (st *SwipeableTabs) Dragged(e *fyne.DragEvent) {
	if !st.isDragging {
		st.isDragging = true
		st.startPos = e.Position
		st.totalDX = 0
		st.totalDY = 0
	}
	st.totalDX += e.Dragged.DX
	st.totalDY += e.Dragged.DY
}

// DragEnd evaluates the gesture and triggers tab transition
func (st *SwipeableTabs) DragEnd() {
	if !st.isDragging {
		return
	}
	defer func() {
		st.isDragging = false
		st.totalDX = 0
		st.totalDY = 0
	}()

	currentIdx := st.Tabs.SelectedIndex()
	tabCount := len(st.Tabs.Items)
	width := st.Size().Width

	newIdx, transitioned := EvaluateSwipeTransition(
		st.startPos.X, st.startPos.Y, width,
		st.totalDX, st.totalDY,
		currentIdx, tabCount,
	)

	if transitioned && newIdx != currentIdx {
		st.Tabs.SelectIndex(newIdx)
		if st.OnTabChange != nil {
			st.OnTabChange(newIdx)
		}
	}
}
