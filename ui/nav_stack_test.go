package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestNavigationManager_StackOperations(t *testing.T) {
	test.NewApp()
	nm := ResetNavManagerForTest(nil)

	if nm.CanPop() {
		t.Fatalf("expected CanPop() to be false on empty stack")
	}
	if nm.Depth() != 0 {
		t.Fatalf("expected Depth() 0, got %d", nm.Depth())
	}

	// Underflow protection: Pop on empty stack should return false
	if nm.Pop() {
		t.Fatalf("expected Pop() on empty stack to return false")
	}

	// Push level 1
	backCalled1 := false
	nm.Push(NavEntry{
		Title:       "Idea 1",
		Breadcrumbs: []string{"💡 Ideas", "Idea 1"},
		View:        widget.NewLabel("Idea 1 Content"),
		OnBack:      func() { backCalled1 = true },
	})

	if !nm.CanPop() {
		t.Fatalf("expected CanPop() to be true after Push")
	}
	if nm.Depth() != 1 {
		t.Fatalf("expected Depth() 1, got %d", nm.Depth())
	}

	// Push sub-level 2
	backCalled2 := false
	nm.Push(NavEntry{
		Title:       "Sub-Thought 1",
		Breadcrumbs: []string{"💡 Ideas", "Idea 1", "Sub-Thought 1"},
		View:        widget.NewLabel("Sub-Thought 1 Content"),
		OnBack:      func() { backCalled2 = true },
	})
	if nm.Depth() != 2 {
		t.Fatalf("expected Depth() 2, got %d", nm.Depth())
	}

	// Push sub-level 3 (deep hierarchy)
	backCalled3 := false
	nm.Push(NavEntry{
		Title:       "Sub-Thought 2",
		Breadcrumbs: []string{"💡 Ideas", "Idea 1", "Sub-Thought 1", "Sub-Thought 2"},
		View:        widget.NewLabel("Sub-Thought 2 Content"),
		OnBack:      func() { backCalled3 = true },
	})
	if nm.Depth() != 3 {
		t.Fatalf("expected Depth() 3, got %d", nm.Depth())
	}

	// Pop level 3
	if !nm.Pop() {
		t.Fatalf("expected Pop() to return true")
	}
	if !backCalled3 {
		t.Fatalf("expected OnBack for level 3 to be called")
	}
	if nm.Depth() != 2 {
		t.Fatalf("expected Depth() 2 after 1 pop, got %d", nm.Depth())
	}

	// ResetToRoot
	nm.ResetToRoot()
	if nm.Depth() != 0 {
		t.Fatalf("expected Depth() 0 after ResetToRoot, got %d", nm.Depth())
	}
	if !backCalled2 || !backCalled1 {
		t.Fatalf("expected OnBack for all remaining levels to be called during ResetToRoot")
	}
	if nm.CanPop() {
		t.Fatalf("expected CanPop() to be false after ResetToRoot")
	}
}

func TestNavigationManager_PopToLevel(t *testing.T) {
	nm := ResetNavManagerForTest(nil)

	// Push 4 levels
	for i := 1; i <= 4; i++ {
		nm.Push(NavEntry{
			Title:       "Level",
			Breadcrumbs: []string{"Root", "L1", "L2", "L3"},
			View:        widget.NewLabel("Content"),
		})
	}
	if nm.Depth() != 4 {
		t.Fatalf("expected depth 4, got %d", nm.Depth())
	}

	// PopToLevel 2
	nm.PopToLevel(2)
	if nm.Depth() != 2 {
		t.Fatalf("expected depth 2, got %d", nm.Depth())
	}

	// Out-of-bounds PopToLevel should be safe no-op
	nm.PopToLevel(10)
	if nm.Depth() != 2 {
		t.Fatalf("expected depth to remain 2, got %d", nm.Depth())
	}
}
