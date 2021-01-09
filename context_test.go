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

	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Background(color.Gray{Y: 220})
			proc.Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			proc.Stroke(color.RGBA{B: 255, A: 255})

			proc.Push()
			proc.Fill(color.RGBA{G: 255, A: 255})
			{
				proc.Push()
				proc.Background(color.Black)
				proc.Pop()
			}
			proc.TextSize(10)
			proc.Stroke(color.RGBA{R: 255, A: 255})
			proc.Rect(20, 20, 160, 160)
			proc.Text("sub-context", 25, 100)
			proc.Pop()

			proc.Rect(120, 120, 70, 70)
			proc.TextSize(15)
			proc.Text("global", 125, 150)
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

	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
			proc.Stroke(color.RGBA{B: 255, A: 255})
			proc.Rect(10, 150, 70, 50)

			for i := 1; i < 10; i++ {
				proc.Push()
				proc.Rotate(float64(i) * math.Pi / 30)
				proc.Stroke(color.RGBA{
					B: uint8((i-1)%2) * 255,
					A: 255,
				})
				proc.Fill(nil)
				proc.Rect(10, 150, 70, 50)
				proc.Pop()
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

	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Background(color.Gray{Y: 220})
		},
		func(proc *Proc) {
			proc.Fill(color.RGBA{B: 255, A: 128})
			proc.Stroke(color.RGBA{B: 255, A: 128})

			proc.Push()
			proc.Scale(0.5, 1)
			proc.Fill(color.RGBA{R: 255, A: 128})
			proc.Stroke(color.RGBA{R: 255, A: 128})
			proc.Rect(30, 20, 50, 50)
			proc.Pop()

			proc.Push()
			proc.Fill(color.RGBA{G: 255, A: 128})
			proc.Stroke(color.RGBA{G: 255, A: 128})
			proc.Scale(0.5, 1.3)
			proc.Rect(30, 20, 50, 50)
			proc.Pop()

			proc.Rect(30, 20, 50, 50)
		},
		"testdata/scale.png",
	)

	proc.Run(t)
}
