// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"

	"github.com/go-p5/p5"
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(400, 400)
	p5.Stroke(nil)
	p5.Background(color.Gray{Y: 220})
}

var (
	xs = make([]float64, 50)
	ys = make([]float64, 50)
)

func draw() {
	for i := 1; i < len(xs); i++ {
		xs[i-1] = xs[i]
		ys[i-1] = ys[i]
	}

	xs[len(xs)-1] = p5.Event.Mouse.Position.X
	ys[len(xs)-1] = p5.Event.Mouse.Position.Y

	for i := range xs {
		p5.Fill(color.RGBA{R: 255, A: uint8(i * 5)})
		p5.Ellipse(xs[i], ys[i], float64(i), float64(i))
	}
}
