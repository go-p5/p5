// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image"
	"image/color"
	"math"
	"reflect"
	"testing"

	"golang.org/x/image/draw"
)

func TestPhysCanvas(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestGProc(t, w, h,
		func(*Proc) {
			PhysCanvas(400, 200, -20, +20, -10, +10)
			Background(color.Gray{Y: 220})
		},
		func(p *Proc) {
			StrokeWidth(2)
			Fill(color.RGBA{R: 255, A: 208})
			Ellipse(0, 0, 20, 10)

			for _, tc := range []struct {
				usr, sys float64
			}{
				{-20, 0},
				{0, 200},
				{+20, 400},
			} {
				if got, want := p.cfg.u2sX(tc.usr), tc.sys; got != want {
					t.Errorf(
						"invalid usr->sys X-conversion: got=%v, want=%v",
						got, want,
					)
				}
				if got, want := p.cfg.s2uX(tc.sys), tc.usr; got != want {
					t.Errorf(
						"invalid sys->usr X-conversion: got=%v, want=%v",
						got, want,
					)
				}
			}

			for _, tc := range []struct {
				usr, sys float64
			}{
				{-10, 0},
				{0, 100},
				{+10, 200},
			} {
				if got, want := p.cfg.u2sY(tc.usr), tc.sys; got != want {
					t.Errorf(
						"invalid usr->sys Y-conversion: got=%v, want=%v",
						got, want,
					)
				}
				if got, want := p.cfg.s2uY(tc.sys), tc.usr; got != want {
					t.Errorf(
						"invalid sys->usr Y-conversion: got=%v, want=%v",
						got, want,
					)
				}
			}

		},
		"testdata/canvas.png",
	)
	proc.Run(t)
}

func TestAPIShapes(t *testing.T) {
	old := gproc
	defer func() {
		gproc = old
	}()

	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(*Proc) {
			Canvas(400, 400)
			Background(color.Gray{Y: 220})
		},
		func(*Proc) {
			StrokeWidth(2)
			Fill(color.RGBA{R: 255, A: 208})
			Ellipse(50, 50, 80, 80)

			Fill(color.RGBA{B: 255, A: 208})
			Circle(50, 50, 40)

			Fill(color.RGBA{B: 255, A: 208})
			Quad(50, 50, 80, 50, 80, 120, 60, 120)

			Fill(color.RGBA{G: 255, A: 208})
			Rect(200, 200, 50, 100)

			Fill(color.RGBA{B: 255, A: 208})
			Square(220, 130, 40)

			Fill(color.RGBA{G: 255, A: 208})
			Triangle(100, 100, 120, 120, 80, 120)

			TextSize(24)
			Text("Hello, World!", 10, 300)

			Stroke(color.Black)
			StrokeWidth(5)
			Arc(300, 100, 80, 20, 0, 1.5*math.Pi)

			Stroke(color.RGBA{R: 255, A: 128})
			Line(300, 0, 300, 400)
		},
		"testdata/api_shapes.png",
	)
	gproc = proc.Proc
	proc.Run(t)
}

func TestBezier(t *testing.T) {
	old := gproc
	defer func() {
		gproc = old
	}()

	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(*Proc) {
			Canvas(400, 400)
			Background(color.Gray{Y: 220})
		},
		func(*Proc) {
			Stroke(color.RGBA{B: 255, A: 128})
			StrokeWidth(3)
			Bezier(100, 340, 230, 360, 100, 80, 300, 260)

			Stroke(color.RGBA{R: 255, A: 208})
			StrokeWidth(5)
			Bezier(100, 100, 230, 80, 100, 30, 300, 200)
		},
		"testdata/api_shapes_bezier.png",
	)
	gproc = proc.Proc
	proc.Run(t)
}

func TestCurve(t *testing.T) {
	const (
		w = 400
		h = 400
	)
	proc := newTestGProc(t, w, h,
		func(*Proc) {
			Background(color.Gray{Y: 220})
		},
		func(*Proc) {
			CurveTightness(0)
			Stroke(color.RGBA{R: 255, A: 255})
			Curve(5, 26, 5, 26, 73, 24, 73, 61)

			Push()
			Stroke(color.Black)
			Curve(5, 26, 73, 24, 73, 61, 15, 65)
			Pop()

			Curve(73, 24, 73, 61, 15, 65, 15, 65)
			Text("tau=0", 10, 50)

			{
				Push()
				Translate(100, 100)
				CurveTightness(1)
				Stroke(color.RGBA{R: 255, A: 255})
				Curve(5, 26, 5, 26, 73, 24, 73, 61)

				Push()
				{
					Stroke(color.Black)
					Curve(5, 26, 73, 24, 73, 61, 15, 65)
				}
				Pop()

				Curve(73, 24, 73, 61, 15, 65, 15, 65)
				Text("tau=+1", 10, 50)

				Pop()
			}

			{
				Push()
				Translate(200, 200)
				CurveTightness(-2)
				Stroke(color.RGBA{R: 255, A: 255})
				Curve(5, 26, 5, 26, 73, 24, 73, 61)

				Push()
				{
					Stroke(color.Black)
					Curve(5, 26, 73, 24, 73, 61, 15, 65)
				}
				Pop()

				Curve(73, 24, 73, 61, 15, 65, 15, 65)
				Text("tau=-2", 10, 50)

				Pop()
			}
			{
				Push()
				Translate(300, 300)
				CurveTightness(+2)
				Stroke(color.RGBA{R: 255, A: 255})
				Curve(5, 26, 5, 26, 73, 24, 73, 61)

				Push()
				{
					Stroke(color.Black)
					Curve(5, 26, 73, 24, 73, 61, 15, 65)
				}
				Pop()

				Curve(73, 24, 73, 61, 15, 65, 15, 65)
				Text("tau=+2", 10, 50)

				Pop()
			}
		},
		"testdata/api_shapes_curve.png",
	)
	proc.Run(t)
}

func TestRandom(t *testing.T) {
	proc := newProc(100, 100)
	proc.RandomSeed(1)

	tests := []struct {
		min  float64
		max  float64
		want float64
	}{
		{0, 5, 1.4393455630422243},
		{-5, -0, -0.6158681223759013},
		{1, 4, 3.3935187588071685},
		{-4, -1, -3.0145149335211046},
		{0, 0, 0},
		{-1, -1, -1},
		{1, 1, 1},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := proc.Random(tt.min, tt.max); got != tt.want {
				t.Errorf("Random() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomGaussian(t *testing.T) {
	proc := newProc(100, 100)
	proc.RandomSeed(1)

	tests := []struct {
		mean   float64
		stdDev float64
		want   float64
	}{
		{0, 1, 0.594696832665853},
		{0, 2, 0.11156098724802527},
		{1, 1, 0.9043068613187665},
		{-1, 2, -3.6908077437524214},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := proc.RandomGaussian(tt.mean, tt.stdDev); got != tt.want {
				t.Errorf("RandomGaussian() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomSeed(t *testing.T) {
	// generatedSequences holds 3 different generated sequences of 5 numbers
	generatedSequences := make([][]float64, 3)

	p1 := newProc(100, 100)
	// Generate the first sequence without calling RandomSeed
	generatedSequences[0] = make([]float64, 5)
	for x := 0; x < 5; x++ {
		generatedSequences[0][x] = p1.Random(0, 1)
	}

	// Generate the second sequence after calling RandomSeed with default seed
	p2 := newProc(100, 100)
	p2.RandomSeed(defaultSeed)
	generatedSequences[1] = make([]float64, 5)
	for x := 0; x < 5; x++ {
		generatedSequences[1][x] = p2.Random(0, 1)
	}

	// Generate the third sequence after calling RandomSeed with seed other than default seed
	p3 := newProc(100, 100)
	p3.RandomSeed(defaultSeed + 1)
	generatedSequences[2] = make([]float64, 5)
	for x := 0; x < 5; x++ {
		generatedSequences[2][x] = p3.Random(0, 1)
	}

	// generatedSequences[0] and generatedSequences[1] should be the same
	if !reflect.DeepEqual(generatedSequences[0], generatedSequences[1]) {
		t.Logf("%v %v", generatedSequences[0], generatedSequences[1])
		t.Errorf("Not calling RandomSeed and calling RandomSeed with the default seed, should produce the same sequence of numbers")
	}
	// generatedSequences[1] and generatedSequences[2] should be different
	if reflect.DeepEqual(generatedSequences[1], generatedSequences[2]) {
		t.Errorf("Calling RandomSeed with different seeds should produce different sequence of numbers")
	}
}

func TestReadImage(t *testing.T) {
	for _, fname := range []string{
		"testdata/gopher.png",
		"testdata/gopher.bmp",
		"testdata/gopher.jpg",
		"testdata/gopher.gif",
		"testdata/gopher.tiff",
	} {
		t.Run(fname, func(t *testing.T) {
			_, err := ReadImage(fname)
			if err != nil {
				t.Fatalf("could not read image %q: %+v", fname, err)
			}
		})
	}
}

func TestDrawImage(t *testing.T) {
	src, err := ReadImage("testdata/gopher.png")
	if err != nil {
		t.Fatalf("could not read image: %+v", err)
	}

	const (
		w = 200
		h = 200
	)

	dst := image.NewRGBA(image.Rect(0, 0, w/2, h/2))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	proc := newTestGProc(t, w, h,
		func(*Proc) {
			Background(color.Gray{Y: 220})
		},
		func(p *Proc) {
			DrawImage(dst, w/2, h/2)

			Push()
			Rotate(-math.Pi / 8)
			Scale(0.6, 0.6)
			Translate(w/2, h*2/3)
			DrawImage(dst, 0, 0)
			Pop()

			for i := 0; i < 10; i++ {
				f := 0.1 * float64(i+1)
				Push()
				Scale(f, f)
				DrawImage(dst, f*float64(w/2), f*float64(h/2))
				Pop()
			}
		},
		"testdata/gopher-image.png",
	)
	proc.Run(t)
}
