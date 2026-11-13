package ui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
)

func TestGraphicPrimitivePool(t *testing.T) {
	pool := NewGraphicPrimitivePool()

	// 1. Circle Pool
	c1 := pool.AcquireCircle(color.White, color.Black, 2.0)
	if c1 == nil {
		t.Fatalf("expected non-nil circle")
	}
	if c1.Hidden {
		t.Fatalf("expected acquired circle not to be hidden")
	}
	pool.ReleaseCircle(c1)
	if !c1.Hidden {
		t.Fatalf("expected released circle to be hidden")
	}

	c2 := pool.AcquireCircle(color.RGBA{R: 255, A: 255}, color.Transparent, 1.0)
	if c2 == nil {
		t.Fatalf("expected non-nil circle from pool")
	}
	pool.ReleaseCircle(c2)

	// 2. Line Pool
	l1 := pool.AcquireLine(color.White, 1.5, fyne.NewPos(0, 0), fyne.NewPos(100, 100))
	if l1 == nil {
		t.Fatalf("expected non-nil line")
	}
	if l1.Position2.X != 100 {
		t.Fatalf("expected p2.X to be 100, got %f", l1.Position2.X)
	}
	pool.ReleaseLine(l1)
	if !l1.Hidden {
		t.Fatalf("expected released line to be hidden")
	}

	// 3. Text Pool
	txt := pool.AcquireText("Node Label", color.White, 14.0, true)
	if txt == nil || txt.Text != "Node Label" {
		t.Fatalf("expected valid text from pool")
	}
	pool.ReleaseText(txt)
	if !txt.Hidden {
		t.Fatalf("expected released text to be hidden")
	}
}

