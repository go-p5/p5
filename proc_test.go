// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"flag"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"testing"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/go-p5/p5/internal/cmpimg"
)

var GenerateTestData = flag.Bool("regen", false, "Uses the current state to regenerate the test data.")

type testWindow struct {
	evts chan event.Event
	opts []app.Option
}

func (w testWindow) Events() <-chan event.Event {
	return w.evts
}

func (w testWindow) Close() {
	w.evts <- system.DestroyEvent{}
}

func (w testWindow) Invalidate() {
	w.evts <- system.FrameEvent{
		Frame: func(ops *op.Ops) {},
	}
}

type testProc struct {
	*Proc
	global bool
	w, h   int
	fname  string
	evts   chan event.Event
}

func newTestGProc(t *testing.T, w, h int, setup, draw func(p *Proc), fname string) *testProc {
	p := newTestProc(t, w, h, setup, draw, fname)
	p.global = true
	return p
}

func newTestProc(t *testing.T, w, h int, setup, draw func(p *Proc), fname string) *testProc {
	t.Helper()

	evts := make(chan event.Event)
	p := newProc(w, h)
	p.Setup = func() { setup(p) }
	p.Draw = func() { draw(p) }
	p.newWindow = func(opts ...app.Option) gioWindow {
		return testWindow{
			evts: evts,
			opts: opts,
		}
	}

	return &testProc{
		Proc:  p,
		w:     w,
		h:     h,
		fname: fname,
		evts:  evts,
	}
}

func (p *testProc) Run(t *testing.T, evts ...event.Event) {
	t.Helper()

	if p.global {
		old := gproc
		defer func() {
			gproc = old
		}()
		gproc = p.Proc
	}

	errc := make(chan error, 1)
	done := make(chan int)
	quit := make(chan int)
	defer close(quit)

	go func() {
		defer close(done)
		select {
		case errc <- p.Proc.run():
		case <-quit:
		}
	}()

	cmds := make([]event.Event, len(evts), len(evts)+2)
	copy(cmds, evts)
	cmds = append(cmds,
		p.frame(t, nil),
		system.DestroyEvent{},
	)

loop:
	for _, evt := range cmds {
		select {
		case p.evts <- evt:
		case <-done:
			// slice of user provided events closed the run-loop.
			break loop
		}
	}

	err := <-errc
	if err != nil {
		t.Fatalf("could not properly shutdown proc: %#+v", err)
	}
}

func (p *testProc) frame(t *testing.T, frame func(ops *op.Ops)) event.Event {
	if frame == nil {
		frame = func(ops *op.Ops) {
			p.screenshot(t)
		}
	}

	return system.FrameEvent{
		Size:  image.Point{X: p.w, Y: p.h},
		Frame: frame,
	}
}

func (p *testProc) screenshot(t *testing.T) {
	if p.fname == "" {
		return
	}

	err := p.Proc.Screenshot(p.fname)
	if err != nil {
		t.Fatalf("could not take screenshot: %+v", err)
	}

	fname := p.fname
	got, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("could not read back screenshot: %+v", err)
	}

	ext := filepath.Ext(fname)
	fname = fname[:len(fname)-len(ext)] + "_golden" + ext

	if *GenerateTestData {
		err = ioutil.WriteFile(fname, got, 0644)
		if err != nil {
			t.Fatalf("could not regen reference file %q: %+v", fname, err)
		}
	}

	want, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("could not read back golden: %+v", err)
	}

	ok, err := cmpimg.EqualApprox(ext[1:], got, want, 0.1)
	if err != nil {
		t.Fatalf("%s: could not compare images: %+v", p.fname, err)
	}
	if !ok {
		t.Fatalf("%s: images compare different", p.fname)
	}

	if err := os.Remove(p.fname); err != nil {
		t.Logf("could not delete image %s, err: %s", p.fname, err)
	}

}

func TestText(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Background(color.Gray{Y: 220})
			proc.Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			proc.Rect(20, 20, 160, 160)
			proc.TextSize(25)
			proc.Text("Hello, World!", 25, 100)
		},
		"testdata/text.png",
	)

	proc.Run(t)
}

func TestHelloWorld(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(p5 *Proc) {
			p5.Canvas(400, 400)
			p5.Background(color.Gray{Y: 220})
		},
		func(p5 *Proc) {
			p5.StrokeWidth(2)
			p5.Fill(color.RGBA{R: 255, A: 208})
			p5.Ellipse(50, 50, 80, 80)

			p5.Fill(color.RGBA{B: 255, A: 208})
			p5.Quad(50, 50, 80, 50, 80, 120, 60, 120)

			p5.Fill(color.RGBA{G: 255, A: 208})
			p5.Rect(200, 200, 50, 100)

			p5.Fill(color.RGBA{G: 255, A: 208})
			p5.Triangle(100, 100, 120, 120, 80, 120)

			p5.TextSize(24)
			p5.Text("Hello, World!", 10, 300)

			p5.Stroke(color.Black)
			p5.StrokeWidth(5)
			p5.Arc(300, 100, 80, 20, 0, 1.5*math.Pi)
		},
		"testdata/hello.png",
	)
	proc.Run(t)
}

func TestShutdown(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(*Proc) {},
		func(*Proc) {},
		"",
	)
	proc.Run(t,
		proc.frame(t, nil), proc.frame(t, nil),
		system.DestroyEvent{},
		proc.frame(t, func(ops *op.Ops) {
			t.Fatalf("should not have executed this frame")
		}),
	)
}
