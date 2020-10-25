// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Ellipse draws an ellipse at (x,y) with the provided width and height.
func Ellipse(x, y, w, h float64) {
	defer op.Push(gctx.ctx.Ops).Pop()

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
	path.End().Add(gctx.ctx.Ops)

	r32 := gctx.rect()
	clr := gctx.cfg.fill

	paint.ColorOp{Color: rgba(clr)}.Add(gctx.ctx.Ops)
	paint.PaintOp{Rect: r32}.Add(gctx.ctx.Ops)
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

// Quad draws a quadrilateral.
func Quad(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	defer op.Push(gctx.ctx.Ops).Pop()

	var (
		p1 = pt32(x1, y1)
		p2 = pt32(x2, y2)
		p3 = pt32(x3, y3)
		p4 = pt32(x4, y4)
	)

	var path clip.Path
	path.Begin(gctx.ctx.Ops)
	path.Move(p1)
	path.Line(p2.Sub(p1))
	path.Line(p3.Sub(p2))
	path.Line(p4.Sub(p3))
	path.Line(p1.Sub(p4))
	path.End().Add(gctx.ctx.Ops)

	r32 := gctx.rect()
	clr := gctx.cfg.fill

	paint.ColorOp{Color: rgba(clr)}.Add(gctx.ctx.Ops)
	paint.PaintOp{Rect: r32}.Add(gctx.ctx.Ops)
}
