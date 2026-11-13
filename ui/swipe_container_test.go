package ui

import (
	"testing"
)

func TestEvaluateSwipeTransition_StandardTabs(t *testing.T) {
	tabCount := 4 // Ideas (0), Memories (1), Mind Map (2), Settings (3)
	width := float32(400)

	// 1. Below threshold -> No transition
	newTab, switched := EvaluateSwipeTransition(100, 200, width, -20, 0, 0, tabCount)
	if switched || newTab != 0 {
		t.Fatalf("expected no switch below threshold (dx=-20), got tab=%d switched=%v", newTab, switched)
	}

	// 2. Vertical drag dominance -> No transition
	newTab, switched = EvaluateSwipeTransition(100, 200, width, -50, 100, 0, tabCount)
	if switched || newTab != 0 {
		t.Fatalf("expected no switch for vertical dominance, got tab=%d switched=%v", newTab, switched)
	}

	// 3. Valid Left Swipe on Ideas (0 -> 1)
	newTab, switched = EvaluateSwipeTransition(100, 200, width, -60, 5, 0, tabCount)
	if !switched || newTab != 1 {
		t.Fatalf("expected switch from 0 to 1 on swipe left, got tab=%d switched=%v", newTab, switched)
	}

	// 4. Valid Right Swipe on Memories (1 -> 0)
	newTab, switched = EvaluateSwipeTransition(100, 200, width, 60, 5, 1, tabCount)
	if !switched || newTab != 0 {
		t.Fatalf("expected switch from 1 to 0 on swipe right, got tab=%d switched=%v", newTab, switched)
	}

	// 5. Left boundary underflow (Swiping right at tab 0)
	newTab, switched = EvaluateSwipeTransition(100, 200, width, 60, 0, 0, tabCount)
	if switched || newTab != 0 {
		t.Fatalf("expected underflow clamp at tab 0, got tab=%d switched=%v", newTab, switched)
	}

	// 6. Right boundary overflow (Swiping left at tab 3)
	newTab, switched = EvaluateSwipeTransition(100, 200, width, -60, 0, 3, tabCount)
	if switched || newTab != 3 {
		t.Fatalf("expected overflow clamp at tab 3, got tab=%d switched=%v", newTab, switched)
	}
}

func TestEvaluateSwipeTransition_MindMapConflictGuard(t *testing.T) {
	tabCount := 4
	width := float32(400)
	mindMapTab := 2

	// 1. Center canvas drag -> MUST BE BLOCKED (Map panning reserved)
	// startX = 200 (center), startY = 300 (center)
	newTab, switched := EvaluateSwipeTransition(200, 300, width, -80, 0, mindMapTab, tabCount)
	if switched || newTab != mindMapTab {
		t.Fatalf("expected center drag on Mind Map to be blocked for canvas panning, got tab=%d switched=%v", newTab, switched)
	}

	// 2. Left Edge Swipe (< 24dp) -> ALLOWED (Swiping right from left edge: 2 -> 1)
	// startX = 15 (< 24)
	newTab, switched = EvaluateSwipeTransition(15, 300, width, 80, 0, mindMapTab, tabCount)
	if !switched || newTab != 1 {
		t.Fatalf("expected left edge swipe on Mind Map to switch to tab 1, got tab=%d switched=%v", newTab, switched)
	}

	// 3. Right Edge Swipe (> width - 24dp) -> ALLOWED (Swiping left from right edge: 2 -> 3)
	// startX = 390 (> 400 - 24 = 376)
	newTab, switched = EvaluateSwipeTransition(390, 300, width, -80, 0, mindMapTab, tabCount)
	if !switched || newTab != 3 {
		t.Fatalf("expected right edge swipe on Mind Map to switch to tab 3, got tab=%d switched=%v", newTab, switched)
	}

	// 4. Header Swipe (startY < 55dp) -> ALLOWED (Swiping on top header)
	newTab, switched = EvaluateSwipeTransition(200, 30, width, -80, 0, mindMapTab, tabCount)
	if !switched || newTab != 3 {
		t.Fatalf("expected header swipe on Mind Map to switch to tab 3, got tab=%d switched=%v", newTab, switched)
	}
}

