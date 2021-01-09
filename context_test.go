// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
	"math"
	"testing"
)

func TestPushPop(t *testing.T) {
	const (
		w = 200
		h = 200
	)

	proc := newTestGProc(t, w, h,
		func(proc *Proc) {
			Background(color.Gray{Y: 220})
			Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			Stroke(color.RGBA{B: 255, A: 255})

			Push()
			Fill(color.RGBA{G: 255, A: 255})
			{
				Push()
				Background(color.Black)
				Pop()
			}
			TextSize(10)
			Stroke(color.RGBA{R: 255, A: 255})
			Rect(20, 20, 160, 160)
			Text("sub-context", 25, 100)
			Pop()

			Rect(120, 120, 70, 70)
			TextSize(15)
			Text("global", 125, 150)
		},
		"testdata/push-pop.png",
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
	)

	proc.Run(t)
}
