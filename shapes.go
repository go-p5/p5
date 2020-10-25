// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Ellipse draws an ellipse at (x,y) with the provided width and height.
func Ellipse(x, y, w, h float64) {
	w *= 0.5
	h *= 0.5

	var (
		ec float64
		f1 f32.Point
		f2 f32.Point
		p1 = pt32(x-w, y)
	)

	switch {
	case math.Abs(w) > math.Abs(h):
		ec = math.Sqrt(w*w - h*h)
		f1 = pt32(x+ec, y).Sub(p1)
		f2 = pt32(x-ec, y).Sub(p1)
	default:
		ec = math.Sqrt(h*h - w*w)
		f1 = pt32(x, y+ec).Sub(p1)
		f2 = pt32(x, y-ec).Sub(p1)
	}

	var path clip.Path
	path.Begin(gctx.ctx.Ops)
	path.Move(p1)
	path.Arc(f1, f2, 2*math.Pi)

	var (
		shape = path.End()
		clr   = gctx.cfg.fill
	)
	paint.FillShape(gctx.ctx.Ops, shape, rgba(clr))
}

// Circle draws a circle at (x,y) with a diameter d.
func Circle(x, y, d float64) {
	Ellipse(x, y, d, d)
}

// Arc draws an ellipsoidal arc centered at (x,y), with the provided
// width and height, and a path from the beg to end radians.
// Positive angles denote a counter-clockwise path.
func Arc(x, y, w, h float64, beg, end float64) {
	panic("not implemented")
}

// Line draws a line between (x1,y1) and (x2,y2).
func Line(x1, y1, x2, y2 float64) {
	panic("not implemented")
}

// Quad draws a quadrilateral, connecting the 4 points (x1,y1),
// (x2,y2), (x3,y3) and (x4,y4) together.
func Quad(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	poly(
		pt32(x1, y1),
		pt32(x2, y2),
		pt32(x3, y3),
		pt32(x4, y4),
	)
}

// Rect draws a rectangle at (x,y) with width w and height h.
func Rect(x, y, w, h float64) {
	Quad(x, y, x+w, y, x+w, y+h, x, y+h)
}

// Square draws a square at (x,y) with size s.
func Square(x, y, s float64) {
	Rect(x, y, s, s)
}

// Triangle draws a triangle, connecting the 3 points (x1,y1), (x2,y2)
// and (x3,y3) together.
func Triangle(x1, y1, x2, y2, x3, y3 float64) {
	poly(
		pt32(x1, y1),
		pt32(x2, y2),
		pt32(x3, y3),
	)
}

func poly(ps ...f32.Point) {
	var path clip.Path
	path.Begin(gctx.ctx.Ops)
	path.Move(ps[0])
	for _, p := range ps[1:] {
		path.Line(p.Sub(path.Pos()))
	}

	var (
		poly = path.End()
		clr  = rgba(gctx.cfg.fill)
	)

	paint.FillShape(gctx.ctx.Ops, poly, clr)
}
