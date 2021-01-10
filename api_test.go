// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
	"math"
	"testing"
)

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

func TestRandom(t *testing.T) {
	old := gproc
	defer func() {
		gproc = old
	}()

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
			if got := Random(tt.min, tt.max); got != tt.want {
				t.Errorf("Random() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomGaussian(t *testing.T) {
	old := gproc
	defer func() {
		gproc = old
	}()

	tests := []struct {
		mean   float64
		stdDev float64
		want   float64
	}{
		{0, 1, 0.06817116484613528},
		{0, 2, -1.317039245855091},
		{1, 1, 0.8714616847016733},
		{-1, 2, 0.20863494591954312},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := RandomGaussian(tt.mean, tt.stdDev); got != tt.want {
				t.Errorf("RandomGaussian() = %v, want %v", got, tt.want)
			}
		})
	}
}
