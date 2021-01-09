// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
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
