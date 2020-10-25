// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
	"math"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type gcontext struct {
	w    *app.Window
	size struct {
		x, y int
	}
	cfg struct {
		th   *material.Theme
		bkg  color.Color
		fill color.Color
	}

	ctx layout.Context
}

func newContext(w, h int) *gcontext {
	var gctx gcontext
	gctx.size.x = w
	gctx.size.y = h
	gctx.cfg.bkg = color.Transparent
	gctx.cfg.fill = color.White
	gctx.ctx = layout.Context{
		Ops: new(op.Ops),
	}
	return &gctx
}

func (gctx *gcontext) run(setup, draw Func) error {
	setup()

	if gctx.cfg.th == nil {
		gctx.cfg.th = material.NewTheme(gofont.Collection())
	}

	gctx.w = app.NewWindow(app.Title("p5"), app.Size(
		unit.Px(float32(gctx.size.x)),
		unit.Px(float32(gctx.size.y)),
	))

	go func() {
		tck := time.NewTicker(10 * time.Millisecond)
		defer tck.Stop()
		for range tck.C {
			gctx.w.Invalidate()
		}
	}()

	for {
		e := <-gctx.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				gctx.w.Invalidate()
				gctx.w.Close()
			}

		case pointer.Event:
			switch e.Type {
			case pointer.Press:
				Event.Mouse.Pressed = true
			case pointer.Release:
				Event.Mouse.Pressed = false
			case pointer.Move:
				Event.Mouse.PrevPosition = Event.Mouse.Position
				Event.Mouse.Position.X = float64(e.Position.X)
				Event.Mouse.Position.Y = float64(e.Position.Y)
			}

		case system.FrameEvent:
			gctx.draw(e, draw)
		}
	}
}

func (gctx *gcontext) draw(e system.FrameEvent, draw Func) {
	gctx.ctx = layout.NewContext(gctx.ctx.Ops, e)

	ops := gctx.ctx.Ops
	r32 := gctx.rect()
	clr := rgba(gctx.cfg.bkg)
	paint.ColorOp{Color: clr}.Add(ops)
	paint.PaintOp{Rect: r32}.Add(ops)

	draw()

	e.Frame(gctx.ctx.Ops)
}

func (gctx *gcontext) rect() f32.Rectangle {
	return f32.Rect(
		0, 0,
		float32(gctx.size.x),
		float32(gctx.size.y),
	)
}

func rgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func pt32(x, y float64) f32.Point {
	return f32.Pt(float32(x), float32(y))
}

// Canvas defines the dimensions of the painting area, in pixels.
func Canvas(w, h int) {
	gctx.size.x = w
	gctx.size.y = h
}

// Background defines the background color for the painting area.
// The default color is transparent.
func Background(c color.Color) {
	gctx.cfg.bkg = c
}

// Fill sets the color used to fill shapes.
func Fill(c color.Color) {
	gctx.cfg.fill = c
}

// Ellipse draws an ellipse at (x,y) with the provided width and height.
func Ellipse(x, y, w, h float64) {
	defer op.Push(gctx.ctx.Ops).Pop()

	w *= 0.5
	h *= 0.5

	var (
		ec float64
		f1 f32.Point
		f2 f32.Point
		p1 = pt32(x-w, y)
	)

	switch {
	case math.Abs(w) > math.Abs(h):
		ec = math.Sqrt(w*w - h*h)
		f1 = pt32(x+ec, y).Sub(p1)
		f2 = pt32(x-ec, y).Sub(p1)
	default:
		ec = math.Sqrt(h*h - w*w)
		f1 = pt32(x, y+ec).Sub(p1)
		f2 = pt32(x, y-ec).Sub(p1)
	}

	var path clip.Path
	path.Begin(gctx.ctx.Ops)
	path.Move(p1)
	path.Arc(f1, f2, 2*math.Pi)
	path.End().Add(gctx.ctx.Ops)

	r32 := gctx.rect()
	clr := gctx.cfg.fill

	paint.ColorOp{Color: rgba(clr)}.Add(gctx.ctx.Ops)
	paint.PaintOp{Rect: r32}.Add(gctx.ctx.Ops)

}

// Circle draws a circle at (x,y) with a diameter d.
func Circle(x, y, d float64) {
	Ellipse(x, y, d, d)
}

// Arc draws an ellipsoidal arc centered at (x,y), with the provided
// width and height, and a path from the beg to end radians.
// Positive angles denote a counter-clockwise path.
func Arc(x, y, w, h float64, beg, end float64) {
	panic("not implemented")
}

// Line draws a line between (x1,y1) and (x2,y2).
func Line(x1, y1, x2, y2 float64) {
	panic("not implemented")
}
