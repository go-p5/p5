// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/spatial/r1"
)

const (
	defaultWidth  = 400
	defaultHeight = 400

	defaultFrameRate = 15 * time.Millisecond

	defaultSeed = 1
)

var (
	defaultBkgColor    = color.Transparent
	defaultFillColor   = color.White
	defaultStrokeColor = color.Black

	defaultTextColor = color.Black
	defaultTextSize  = float32(12)
)

// gioWindow represents an operating system window operated by Gio.
type gioWindow interface {
	// Close the window. The window's event loop should exit when it receives
	// system.DestroyEvent.
	Close()

	// Events returns the channel where events are delivered.
	Events() <-chan event.Event

	// Invalidate the window such that a FrameEvent will be generated immediately.
	// If the window is inactive, the event is sent when the window becomes active.
	Invalidate()
}

var _ gioWindow = (*app.Window)(nil)

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
		run  bool
		loop bool
	}
	cfg struct {
		w int
		h int

		x    r1.Interval
		y    r1.Interval
		u2sX func(v float64) float64 // translate from user- to system coords
		u2sY func(v float64) float64 // translate from user- to system coords
		s2uX func(v float64) float64 // translate from system- to user coords
		s2uY func(v float64) float64 // translate from system- to user coords

		th *material.Theme
	}

	ctx  layout.Context
	stk  *stackOps
	head *headless.Window
	rand *rand.Rand

	newWindow func(opts ...app.Option) gioWindow
}

func newProc(w, h int) *Proc {
	proc := &Proc{
		ctx: layout.Context{
			Ops: new(op.Ops),
			Constraints: layout.Constraints{
				Max: image.Pt(w, h),
			},
		},
		rand: rand.New(rand.NewSource(defaultSeed)),

		newWindow: func(opts ...app.Option) gioWindow {
			return app.NewWindow(opts...)
		},
	}
	proc.ctl.FrameRate = defaultFrameRate
	proc.ctl.loop = true
	proc.stk = newStackOps(proc.ctx.Ops)
	proc.initCanvas(w, h)

	proc.cfg.th = material.NewTheme(gofont.Collection())
	proc.stk.cur().stroke.style.Width = 2

	return proc
}

func (p *Proc) initCanvas(w, h int) {
	p.initCanvasDim(w, h, 0, float64(w), 0, float64(h))
	p.stk.cur().bkg = defaultBkgColor
	p.stk.cur().fill = defaultFillColor
	p.stk.cur().stroke.color = defaultStrokeColor

	p.stk.cur().text.color = defaultTextColor
	p.stk.cur().text.align = text.Start
	p.stk.cur().text.size = defaultTextSize
}

func (p *Proc) initCanvasDim(w, h int, xmin, xmax, ymin, ymax float64) {
	p.cfg.w = w
	p.cfg.h = h
	p.cfg.x = r1.Interval{Min: xmin, Max: xmax}
	p.cfg.y = r1.Interval{Min: ymin, Max: ymax}

	var (
		wdx = 1 / (p.cfg.x.Max - p.cfg.x.Min) * float64(w)
		hdy = 1 / (p.cfg.y.Max - p.cfg.y.Min) * float64(h)

		dx = 1 / wdx
		dy = 1 / hdy
	)

	p.cfg.u2sX = func(v float64) float64 {
		return (v - p.cfg.x.Min) * wdx
	}

	p.cfg.s2uX = func(v float64) float64 {
		return (v * dx) + p.cfg.x.Min
	}

	p.cfg.u2sY = func(v float64) float64 {
		return (v - p.cfg.y.Min) * hdy
	}

	p.cfg.s2uY = func(v float64) float64 {
		return (v * dy) + p.cfg.y.Min
	}
}

func (p *Proc) cnvSize() (w, h float64) {
	w = float64(p.cfg.w)
	h = float64(p.cfg.h)
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
	p.setupUserFuncs()

	p.Setup()

	var (
		err    error
		width  = p.cfg.w
		height = p.cfg.h
	)

	w := p.newWindow(app.Title("p5"), app.Size(
		unit.Px(float32(width)),
		unit.Px(float32(height)),
	))
	p.head, err = headless.NewWindow(width, height)
	if err != nil {
		return fmt.Errorf("p5: could not create headless window: %w", err)
	}

	p.ctl.mu.Lock()
	p.ctl.run = true
	p.ctl.mu.Unlock()

	go func() {
		tck := time.NewTicker(p.ctl.FrameRate)
		defer tck.Stop()
		for range tck.C {
			w.Invalidate()
		}
	}()

	var cnt int

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				w.Close()
			case "F11":
				fname := fmt.Sprintf("out-%03d.png", cnt)
				err = p.Screenshot(fname)
				if err != nil {
					log.Printf("could not take screenshot: %+v", err)
				}
				cnt++
			}

		case pointer.Event:
			switch e.Type {
			case pointer.Press:
				Event.Mouse.Pressed = true
			case pointer.Release:
				Event.Mouse.Pressed = false
			case pointer.Move:
				Event.Mouse.PrevPosition = Event.Mouse.Position
				Event.Mouse.Position.X = p.cfg.s2uX(float64(e.Position.X))
				Event.Mouse.Position.Y = p.cfg.s2uY(float64(e.Position.Y))
			}
			Event.Mouse.Buttons = Buttons(e.Buttons)

		case system.FrameEvent:
			if p.loop() {
				p.draw(e)
			}
		}
	}
}

func (p *Proc) setupUserFuncs() {
	if p.Setup == nil {
		p.Setup = func() {}
	}
	if p.Draw == nil {
		p.Draw = func() {}
	}
	if p.Mouse == nil {
		p.Mouse = func() {}
	}
}

func (p *Proc) loop() bool {
	p.ctl.mu.RLock()
	defer p.ctl.mu.RUnlock()
	return p.ctl.loop
}

func (p *Proc) draw(e system.FrameEvent) {
	p.ctx = layout.NewContext(p.ctx.Ops, e)

	ops := p.ctx.Ops
	clr := rgba(p.stk.cur().bkg)
	paint.Fill(ops, clr)

	p.Draw()

	e.Frame(ops)
}

func (p *Proc) pt(x, y float64) f32.Point {
	return f32.Point{
		X: float32(p.cfg.u2sX(x)),
		Y: float32(p.cfg.u2sY(y)),
	}
}

func rgba(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// Canvas defines the dimensions of the painting area, in pixels.
func (p *Proc) Canvas(w, h int) {
	p.initCanvasDim(w, h, 0, float64(w), 0, float64(h))
}

// PhysCanvas sets the dimensions of the painting area, in pixels, and
// associates physical quantities.
func (p *Proc) PhysCanvas(w, h int, xmin, xmax, ymin, ymax float64) {
	p.initCanvasDim(w, h, xmin, xmax, ymin, ymax)
}

// Background defines the background color for the painting area.
// The default color is transparent.
func (p *Proc) Background(c color.Color) {
	p.stk.cur().bkg = c
}

func (p *Proc) doStroke() bool {
	return p.stk.cur().stroke.color != nil &&
		p.stk.cur().stroke.style.Width > 0
}

// Stroke sets the color of the strokes.
func (p *Proc) Stroke(c color.Color) {
	p.stk.cur().stroke.color = c
}

// StrokeWidth sets the size of the strokes.
func (p *Proc) StrokeWidth(v float64) {
	p.stk.cur().stroke.style.Width = float32(v)
}

func (p *Proc) doFill() bool {
	return p.stk.cur().fill != nil
}

// Fill sets the color used to fill shapes.
func (p *Proc) Fill(c color.Color) {
	p.stk.cur().fill = c
}

// TextSize sets the text size.
func (p *Proc) TextSize(size float64) {
	p.stk.cur().text.size = float32(size)
}

// Text draws txt on the screen at (x,y).
func (p *Proc) Text(txt string, x, y float64) {
	x = p.cfg.u2sX(x)
	y = p.cfg.u2sY(y)

	var (
		offset = x
		w, _   = p.cnvSize()
		size   = p.stk.cur().text.size
	)
	switch p.stk.cur().text.align {
	case text.End:
		offset = x - w
	case text.Middle:
		offset = x - 0.5*w
	}
	defer op.Save(p.ctx.Ops).Load()
	op.Offset(f32.Point{
		X: float32(offset),
		Y: float32(y) - size,
	}).Add(p.ctx.Ops) // shift to use baseline

	l := material.Label(p.cfg.th, unit.Px(size), txt)
	l.Color = rgba(p.stk.cur().text.color)
	l.Alignment = p.stk.cur().text.align
	l.Layout(p.ctx)
}

// Screenshot saves the current canvas to the provided file.
// Supported file formats are: PNG, JPEG and GIF.
func (p *Proc) Screenshot(fname string) error {
	err := p.head.Frame(p.ctx.Ops)
	if err != nil {
		return fmt.Errorf("p5: could not run headless frame: %w", err)
	}

	img, err := p.head.Screenshot()
	if err != nil {
		return fmt.Errorf("p5: could not take screenshot: %w", err)
	}

	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("p5: could not create screenshot file: %w", err)
	}
	defer f.Close()

	var encode func(io.Writer, image.Image) error
	switch ext := filepath.Ext(fname); strings.ToLower(ext) {
	case ".jpeg", ".jpg":
		encode = func(w io.Writer, img image.Image) error {
			return jpeg.Encode(w, img, nil)
		}
	case ".gif":
		encode = func(w io.Writer, img image.Image) error {
			return gif.Encode(w, img, nil)
		}
	case ".png":
		encode = png.Encode
	default:
		log.Printf("unknown file extension %q. using png", ext)
		encode = png.Encode
	}

	err = encode(f, img)
	if err != nil {
		return fmt.Errorf("p5: could not encode screenshot: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("p5: could not save screenshot: %w", err)
	}

	return nil
}

// RandomSeed changes the sequence of numbers generated by Random.
func (p *Proc) RandomSeed(seed uint64) {
	p.rand.Seed(seed)
}

// Random returns a pseudo-random number in [min,max).
// Random will produce the same sequence of numbers every time the program runs.
// Use RandomSeed with a seed that changes (like time.Now().UnixNano()) in order to
// produce different sequences of numbers.
func (p *Proc) Random(min, max float64) float64 {
	return p.rand.Float64()*(max-min) + min
}

// RandomGaussian returns a random number following a Gaussian distribution with the provided
// mean and standard deviation.
func (p *Proc) RandomGaussian(mean, stdDev float64) float64 {
	return p.rand.NormFloat64()*stdDev + mean
}
