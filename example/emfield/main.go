// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"

	"github.com/go-p5/p5"
	"gonum.org/v1/gonum/spatial/r2"
)

func main() {
	p5.Run(setup, draw)
}

const (
	width  = 1200
	height = 1200

	xmin = -1.3 * 10e-3
	xmax = +1.3 * 10e-3
	ymin = -1.0 * 10e-3
	ymax = +1.0 * 10e-3

	nlines = 35 // number of lines per particle
)

var (
	ps = []Particle{
		{p: r2.Vec{X: -10e-3 / 2}, q: +1, r: 4 * 10e-5},
		{p: r2.Vec{X: +10e-3 / 2}, q: +1, r: 4 * 10e-5},
	}
	em = newEMField(ps, nlines)
)

func setup() {
	p5.PhysCanvas(width, height, xmin, xmax, ymin, ymax)
	p5.Background(color.Black)
}

func draw() {
	em.update()
	em.draw()
}
