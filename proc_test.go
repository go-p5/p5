// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"
	"testing"

	"gioui.org/app/headless"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/go-p5/p5/internal/cmpimg"
)

type testProc struct {
	*Proc
	w, h  int
	fname string
}

func newTestProc(t *testing.T, w, h int, setup, draw func(p *Proc), fname string) *testProc {
	t.Helper()

	p := newProc(w, h)
	p.Setup = func() { setup(p) }
	p.Draw = func() { draw(p) }

	return &testProc{
		Proc:  p,
		w:     w,
		h:     h,
		fname: fname,
	}
}

func (p *testProc) Run(t *testing.T) {
	t.Helper()

	p.setupUserFuncs()

	p.Proc.Setup()

	var (
		err           error
		width, height = p.cnvSize()
	)

	p.head, err = headless.NewWindow(int(width), int(height))
	if err != nil {
		t.Fatalf("could not create headless window: %+v", err)
	}

	p.Proc.draw(system.FrameEvent{
		Size:  image.Point{X: p.w, Y: p.h},
		Frame: func(ops *op.Ops) {},
	})

	err = p.Proc.Screenshot(p.fname)
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
	want, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("could not read back golden: %+v", err)
	}

	ok, err := cmpimg.EqualApprox(ext[1:], got, want, 0.05)
	if err != nil {
		t.Fatalf("%s: could not compare images: %+v", p.fname, err)
	}
	if !ok {
		t.Fatalf("%s: images compare different", p.fname)
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
