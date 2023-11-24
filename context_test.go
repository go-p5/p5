// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"github.com/go-fonts/latin-modern/lmroman12regular"
)

func TestPushPop(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	fonts := gofont.Collection()

	face, err := opentype.Parse(lmroman12regular.TTF)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %+v", err))
	}
	fonts = append(fonts, font.FontFace{
		Font: font.Font{
			Typeface: "Latin-Modern",
		},
		Face: face,
	})

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			LoadFonts(fonts)
			Background(color.Gray{Y: 220})
			Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			Stroke(color.RGBA{B: 255, A: 255})

			TextSize(15)
			{
				Push()
				Fill(color.RGBA{G: 255, A: 255})
				{
					Push()
					Background(color.Black)
					Fill(color.RGBA{R: 255, A: 255})
					Pop()
				}
				TextFont(font.Font{Typeface: "Latin-Modern"})
				TextSize(20)
				Stroke(color.RGBA{R: 255, A: 255})
				Rect(20, 20, 160, 160)
				Text("sub-context", 25, 100)
				{
					Push()
					TextSize(10)
					Rect(25, 40, 80, 20)
					Text("sub-sub-context", 30, 50)
					Pop()

				}
				Pop()
			}

			Rect(120, 120, 70, 70)
			Text("global", 125, 150)
		},
		"testdata/push-pop.png",
		imgDelta,
	)

	proc.Run(t)
}

func TestRotate(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {
			Fill(color.RGBA{R: 255, A: 255})
			Stroke(color.RGBA{B: 255, A: 255})
			Rect(10, 150, 70, 50)

			for i := 1; i < 10; i++ {
				Push()
				Rotate(float64(i) * math.Pi / 30)
				Stroke(color.RGBA{
					B: uint8((i-1)%2) * 255,
					A: 255,
				})
				Fill(nil)
				Rect(10, 150, 70, 50)
				Pop()
			}
		},
		"testdata/rotate.png",
		imgDelta,
	)

	proc.Run(t)
}

func TestScale(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {
			Fill(color.RGBA{B: 255, A: 128})
			Stroke(color.RGBA{B: 255, A: 128})

			Push()
			Scale(0.5, 1)
			Fill(color.RGBA{R: 255, A: 128})
			Stroke(color.RGBA{R: 255, A: 128})
			Rect(30, 20, 50, 50)
			Pop()

			Push()
			Fill(color.RGBA{G: 255, A: 128})
			Stroke(color.RGBA{G: 255, A: 128})
			Scale(0.5, 1.3)
			Rect(30, 20, 50, 50)
			Pop()

			Rect(30, 20, 50, 50)
		},
		"testdata/scale.png",
		imgDelta,
	)

	proc.Run(t)
}

func TestTranslate(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {
			Fill(color.RGBA{B: 255, A: 128})
			Stroke(color.RGBA{B: 255, A: 128})

			for i := 0; i < 50; i++ {
				Push()
				Translate(float64(i)*5, float64(i)*10)
				Rect(0, 0, 20, 30)
				Pop()
			}
		},
		"testdata/translate.png",
		imgDelta,
	)

	proc.Run(t)
}

func TestShear(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {

			Fill(nil)
			Stroke(color.Black)
			Rect(50, 50, 20, 30)

			Push()
			Fill(color.RGBA{B: 255, A: 255})
			Stroke(color.RGBA{B: 255, A: 255})
			Shear(math.Pi/4, 0)
			Rect(50, 50, 20, 30)
			Text("shear-x", 50, 70)
			Pop()

			Push()
			Fill(color.RGBA{R: 255, A: 255})
			Stroke(color.RGBA{R: 255, A: 255})
			Shear(0, math.Pi/4)
			Rect(50, 50, 20, 30)
			Text("shear-y", 50, 70)
			Pop()
		},
		"testdata/shear.png",
		imgDelta,
	)

	proc.Run(t)
}

func TestMatrix(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {

			Fill(nil)
			Stroke(color.Black)
			Rect(100, 100, 20, 30)

			{
				Push()
				sin, cos := math.Sincos(math.Pi / 6)
				Stroke(color.RGBA{B: 255, A: 255})
				Matrix(cos, sin, -sin, cos, 0, 0)
				Rect(100, 100, 20, 30)
				Pop()
			}

			{
				Push()
				sin, cos := math.Sincos(-math.Pi / 6)
				Stroke(color.RGBA{R: 255, A: 255})
				Matrix(cos, sin, -sin, cos, 0, 0)
				Rect(100, 100, 20, 30)
				Pop()
			}
		},
		"testdata/matrix.png",
		imgDelta,
	)

	proc.Run(t)
}
