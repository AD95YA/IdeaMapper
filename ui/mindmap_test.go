package ui

import (
	"image/color"
	"testing"
)

func TestLevelAccentColor(t *testing.T) {
	// Levels 1 to 6 should return distinct colors
	c1 := LevelAccentColor(1)
	c2 := LevelAccentColor(2)
	c3 := LevelAccentColor(3)
	c4 := LevelAccentColor(4)
	c5 := LevelAccentColor(5)
	c6 := LevelAccentColor(6)

	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	r3, g3, b3, _ := c3.RGBA()

	if r1 == r2 && g1 == g2 && b1 == b2 {
		t.Fatalf("expected level 1 and 2 to have distinct accent colors")
	}
	if r2 == r3 && g2 == g3 && b2 == b3 {
		t.Fatalf("expected level 2 and 3 to have distinct accent colors")
	}

	// Default fallback for invalid level
	defColor := LevelAccentColor(99)
	if defColor != (color.RGBA{R: 148, G: 163, B: 184, A: 255}) {
		t.Fatalf("expected fallback slate color for level 99")
	}
	_ = c4
	_ = c5
	_ = c6
}

func TestCalculateRadialPositions(t *testing.T) {
	// 0 count
	pos0 := CalculateRadialPositions(200, 200, 100, 100, 0)
	if len(pos0) != 0 {
		t.Fatalf("expected empty slice for count 0")
	}

	// 4 nodes arranged around center (200, 200) with radius 100
	pos4 := CalculateRadialPositions(200, 200, 100, 100, 4)
	if len(pos4) != 4 {
		t.Fatalf("expected 4 positions, got %d", len(pos4))
	}

	// Node 0 should be at angle 0 -> (200 + 100, 200) = (300, 200)
	if pos4[0].X < 299 || pos4[0].X > 301 || pos4[0].Y < 199 || pos4[0].Y > 201 {
		t.Fatalf("expected node 0 near (300, 200), got (%f, %f)", pos4[0].X, pos4[0].Y)
	}

	// Node 2 should be at angle Pi -> (200 - 100, 200) = (100, 200)
	if pos4[2].X < 99 || pos4[2].X > 101 || pos4[2].Y < 199 || pos4[2].Y > 201 {
		t.Fatalf("expected node 2 near (100, 200), got (%f, %f)", pos4[2].X, pos4[2].Y)
	}
}

