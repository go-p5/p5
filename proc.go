// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sync"
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
	"gonum.org/v1/gonum/spatial/r1"
)

const (
	defaultWidth  = 400
	defaultHeight = 400

	defaultFrameRate = 15 * time.Millisecond
)

var (
	defaultBkgColor    = color.Transparent
	defaultFillColor   = color.White
	defaultStrokeColor = color.Black
)

// Proc is a p5 processor.
//
// Proc runs the bound Setup function once before the event loop.
// Proc then runs the bound Draw function once per event loop iteration.
type Proc struct {
	Setup Func
	Draw  Func
	Mouse Func

	ctl struct {
		FrameRate time.Duration

		mu   sync.RWMutex
		loop bool
	}
	cfg struct {
		x   r1.Interval
		y   r1.Interval
		trX func(v float64) float64 // translate from user- to system coords
		trY func(v float64) float64 // translate from user- to system coords

		color struct {
			bkg    color.Color
			fill   color.Color
			stroke color.Color
		}

		th    *material.Theme
		linew float32
		sty   clip.StrokeStyle
	}

	ctx layout.Context
}

func newProc(w, h int) *Proc {
	proc := &Proc{
		ctx: layout.Context{
			Ops: new(op.Ops),
			Constraints: layout.Constraints{
				Max: image.Pt(w, h),
			},
		},
	}
	proc.ctl.FrameRate = defaultFrameRate
	proc.ctl.loop = true
	proc.initCanvas(w, h)

	proc.cfg.th = material.NewTheme(gofont.Collection())
	proc.cfg.linew = 2

	return proc
}

func (p *Proc) initCanvas(w, h int) {
	p.cfg.x = r1.Interval{Min: 0, Max: float64(w)}
	p.cfg.y = r1.Interval{Min: 0, Max: float64(h)}
	p.cfg.trX = func(v float64) float64 {
		return (v - p.cfg.x.Min) / (p.cfg.x.Max - p.cfg.x.Min) * float64(w)
	}

	p.cfg.trY = func(v float64) float64 {
		return (v - p.cfg.y.Min) / (p.cfg.y.Max - p.cfg.y.Min) * float64(h)
	}
	p.cfg.color.bkg = defaultBkgColor
	p.cfg.color.fill = defaultFillColor
	p.cfg.color.stroke = defaultStrokeColor
}

func (p *Proc) cnvSize() (w, h float64) {
	w = math.Abs(proc.cfg.x.Max - proc.cfg.x.Min)
	h = math.Abs(proc.cfg.y.Max - proc.cfg.y.Min)
	return w, h
}

func (p *Proc) Run() {
	go func() {
		err := p.run()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (p *Proc) run() error {
	p.Setup()

	width, height := p.cnvSize()

	w := app.NewWindow(app.Title("p5"), app.Size(
		unit.Px(float32(width)),
		unit.Px(float32(height)),
	))

	go func() {
		tck := time.NewTicker(p.ctl.FrameRate)
		defer tck.Stop()
		for range tck.C {
			w.Invalidate()
		}
	}()

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				w.Invalidate()
				w.Close()
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
			if p.loop() {
				p.draw(w, e)
			}
		}
	}
}

func (p *Proc) loop() bool {
	p.ctl.mu.RLock()
	defer p.ctl.mu.RUnlock()
	return p.ctl.loop
}

func (p *Proc) draw(win *app.Window, e system.FrameEvent) {
	p.ctx = layout.NewContext(p.ctx.Ops, e)

	ops := p.ctx.Ops
	clr := rgba(p.cfg.color.bkg)
	paint.Fill(ops, clr)

	p.Draw()

	e.Frame(ops)
}

func (p *Proc) pt(x, y float64) f32.Point {
	return f32.Point{
		X: float32(proc.cfg.trX(x)),
		Y: float32(proc.cfg.trY(y)),
	}
}

func rgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// Canvas defines the dimensions of the painting area, in pixels.
func Canvas(w, h int) {
	proc.initCanvas(w, h)
}

// Background defines the background color for the painting area.
// The default color is transparent.
func Background(c color.Color) {
	proc.cfg.color.bkg = c
}

// Stroke sets the color of the strokes.
func Stroke(c color.Color) {
	proc.cfg.color.stroke = c
}

// Fill sets the color used to fill shapes.
func Fill(c color.Color) {
	proc.cfg.color.fill = c
}
