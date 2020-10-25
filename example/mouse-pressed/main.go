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
	p5.Background(color.Gray{Y: 220})
}

func draw() {
	switch {
	case p5.Event.Mouse.Pressed:
		if p5.Event.Mouse.Buttons.Contain(p5.ButtonLeft) {
			p5.Fill(color.RGBA{R: 255, A: 255})
		}
	default:
		p5.Fill(color.Transparent)
	}
	p5.Ellipse(
		p5.Event.Mouse.Position.X,
		p5.Event.Mouse.Position.Y,
		80, 80,
	)
}
