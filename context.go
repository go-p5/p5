// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
			Event.Mouse.Buttons = Buttons(e.Buttons)

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
