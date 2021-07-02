// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"

	"gioui.org/io/key"
	"github.com/go-p5/p5"
)

func main() {
	p5.Run(setup, draw, p5.WithKeyCallback(kcbk))
}

const (
	width  = 400
	height = 400
)

func setup() {
	p5.Canvas(width, height)
	p5.Stroke(nil)
	p5.Background(color.Gray{Y: 220})
}

var (
	xs = make([]float64, 50)
	ys = make([]float64, 50)

	px, py float64
)

func draw() {
	for i := 1; i < len(xs); i++ {
		xs[i-1] = xs[i]
		ys[i-1] = ys[i]
	}

	xs[len(xs)-1] = px
	ys[len(xs)-1] = py

	for i := range xs {
		p5.Fill(color.RGBA{R: 255, A: uint8(i * 5)})
		p5.Ellipse(xs[i], ys[i], float64(i), float64(i))
	}
}

func kcbk() {
	cur := p5.Event.Key.Cur
	if cur.State != key.Press {
		return
	}
	switch cur.Name {
	case key.NameLeftArrow:
		px -= 5
	case key.NameRightArrow:
		px += 5
	case key.NameUpArrow:
		py -= 5
	case key.NameDownArrow:
		py += 5
	}
	px = clip(px, 0, width)
	py = clip(py, 0, height)
}

func clip(x, min, max float64) float64 {
	switch {
	case x < min:
		x = min
	case x > max:
		x = max
	}
	return x
}
