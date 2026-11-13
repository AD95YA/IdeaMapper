package ui

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// GraphicPrimitivePool manages reusable canvas objects to achieve zero-allocation redraw loops
type GraphicPrimitivePool struct {
	circlePool sync.Pool
	linePool   sync.Pool
	textPool   sync.Pool
}

var globalPrimitivePool = NewGraphicPrimitivePool()

// GetPrimitivePool returns the singleton graphic pool
func GetPrimitivePool() *GraphicPrimitivePool {
	return globalPrimitivePool
}

// NewGraphicPrimitivePool initializes memory pools for canvas primitives
func NewGraphicPrimitivePool() *GraphicPrimitivePool {
	return &GraphicPrimitivePool{
		circlePool: sync.Pool{
			New: func() interface{} {
				c := canvas.NewCircle(color.White)
				return c
			},
		},
		linePool: sync.Pool{
			New: func() interface{} {
				l := canvas.NewLine(color.White)
				return l
			},
		},
		textPool: sync.Pool{
			New: func() interface{} {
				t := canvas.NewText("", color.White)
				return t
			},
		},
	}
}

// AcquireCircle retrieves a recycled Circle from the pool with requested properties
func (p *GraphicPrimitivePool) AcquireCircle(fill color.Color, stroke color.Color, strokeWidth float32) *canvas.Circle {
	c := p.circlePool.Get().(*canvas.Circle)
	c.FillColor = fill
	c.StrokeColor = stroke
	c.StrokeWidth = strokeWidth
	c.Hidden = false
	return c
}

// ReleaseCircle returns a Circle to the pool
func (p *GraphicPrimitivePool) ReleaseCircle(c *canvas.Circle) {
	if c == nil {
		return
	}
	c.Hidden = true
	p.circlePool.Put(c)
}

// AcquireLine retrieves a recycled Line from the pool with requested properties
func (p *GraphicPrimitivePool) AcquireLine(stroke color.Color, strokeWidth float32, p1, p2 fyne.Position) *canvas.Line {
	l := p.linePool.Get().(*canvas.Line)
	l.StrokeColor = stroke
	l.StrokeWidth = strokeWidth
	l.Position1 = p1
	l.Position2 = p2
	l.Hidden = false
	return l
}

// ReleaseLine returns a Line to the pool
func (p *GraphicPrimitivePool) ReleaseLine(l *canvas.Line) {
	if l == nil {
		return
	}
	l.Hidden = true
	p.linePool.Put(l)
}

// AcquireText retrieves a recycled Text object from the pool with requested properties
func (p *GraphicPrimitivePool) AcquireText(text string, col color.Color, textSize float32, bold bool) *canvas.Text {
	t := p.textPool.Get().(*canvas.Text)
	t.Text = text
	t.Color = col
	t.TextSize = textSize
	t.TextStyle = fyne.TextStyle{Bold: bold}
	t.Hidden = false
	return t
}

// ReleaseText returns a Text object to the pool
func (p *GraphicPrimitivePool) ReleaseText(t *canvas.Text) {
	if t == nil {
		return
	}
	t.Hidden = true
	p.textPool.Put(t)
}

