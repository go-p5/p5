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
func (p *Proc) Ellipse(x, y, w, h float64) {
	if !p.doFill() && !p.doStroke() {
		return
	}

	w *= 0.5
	h *= 0.5

	var (
		ec float64
		f1 f32.Point
		f2 f32.Point
		p1 = p.pt(x-w, y)
	)

	switch {
	case math.Abs(w) > math.Abs(h):
		ec = math.Sqrt(w*w - h*h)
		f1 = p.pt(x+ec, y).Sub(p1)
		f2 = p.pt(x-ec, y).Sub(p1)
	default:
		ec = math.Sqrt(h*h - w*w)
		f1 = p.pt(x, y+ec).Sub(p1)
		f2 = p.pt(x, y-ec).Sub(p1)
	}

	path := func(o *op.Ops) clip.PathSpec {
		var path clip.Path
		path.Begin(o)
		path.Move(p1)
		path.Arc(f1, f2, 2*math.Pi)
		return path.End()
	}

	if p.cfg.color.fill != nil {
		stk := op.Push(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.fill),
			clip.Outline{
				Path: path(p.ctx.Ops),
			}.Op(),
		)
		stk.Pop()
	}

	if p.cfg.color.stroke != nil {
		stk := op.Push(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.stroke),
			clip.Stroke{
				Path:  path(p.ctx.Ops),
				Style: p.cfg.stroke,
			}.Op(),
		)
		stk.Pop()
	}
}

// Circle draws a circle at (x,y) with a diameter d.
func (p *Proc) Circle(x, y, d float64) {
	p.Ellipse(x, y, d, d)
}

// Arc draws an ellipsoidal arc centered at (x,y), with the provided
// width and height, and a path from the beg to end radians.
// Positive angles denote a counter-clockwise path.
func (p *Proc) Arc(x, y, w, h float64, beg, end float64) {
	if !p.doStroke() {
		return
	}

	var (
		c  = p.pt(x, y)
		a  = p.cfg.trX(w)
		b  = p.cfg.trY(h)
		f1 f32.Point
		f2 f32.Point
	)

	switch {
	case a >= b:
		f := math.Sqrt(a*a - b*b)
		f1 = c.Add(p.pt(+f, 0))
		f2 = c.Add(p.pt(-f, 0))
	default:
		f := math.Sqrt(b*b - a*a)
		f1 = c.Add(p.pt(0, +f))
		f2 = c.Add(p.pt(0, -f))
	}

	var (
		sin, cos = math.Sincos(beg)
		p0       = p.pt(a*cos, b*sin).Add(c)
		path     clip.Path
	)
	stk := op.Push(p.ctx.Ops)
	path.Begin(p.ctx.Ops)
	path.Move(p0)
	path.Arc(f1.Sub(p0), f2.Sub(p0), float32(end-beg))

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.cfg.color.stroke),
		clip.Stroke{
			Path:  path.End(),
			Style: p.cfg.stroke,
		}.Op(),
	)
	stk.Pop()
}

// Line draws a line between (x1,y1) and (x2,y2).
func (p *Proc) Line(x1, y1, x2, y2 float64) {
	if !p.doStroke() {
		return
	}

	var (
		p1   = p.pt(x1, y1)
		p2   = p.pt(x2, y2)
		path clip.Path
	)
	stk := op.Push(p.ctx.Ops)
	path.Begin(p.ctx.Ops)
	path.Move(p1)
	path.Line(p2.Sub(path.Pos()))

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.cfg.color.stroke),
		clip.Stroke{
			Path:  path.End(),
			Style: p.cfg.stroke,
		}.Op(),
	)
	stk.Pop()
}

// Quad draws a quadrilateral, connecting the 4 points (x1,y1),
// (x2,y2), (x3,y3) and (x4,y4) together.
func (p *Proc) Quad(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	p.poly(
		p.pt(x1, y1),
		p.pt(x2, y2),
		p.pt(x3, y3),
		p.pt(x4, y4),
		p.pt(x1, y1),
	)
}

// Rect draws a rectangle at (x,y) with width w and height h.
func (p *Proc) Rect(x, y, w, h float64) {
	p.Quad(x, y, x+w, y, x+w, y+h, x, y+h)
}

// Square draws a square at (x,y) with size s.
func (p *Proc) Square(x, y, s float64) {
	p.Rect(x, y, s, s)
}

// Triangle draws a triangle, connecting the 3 points (x1,y1), (x2,y2)
// and (x3,y3) together.
func (p *Proc) Triangle(x1, y1, x2, y2, x3, y3 float64) {
	p.poly(
		p.pt(x1, y1),
		p.pt(x2, y2),
		p.pt(x3, y3),
		p.pt(x1, y1),
	)
}

func (p *Proc) Bezier(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	if !p.doStroke() {
		return
	}

	var (
		sp  = p.pt(x1, y1)
		cp0 = p.pt(x2, y2).Sub(sp)
		cp1 = p.pt(x3, y3).Sub(sp)
		ep  = p.pt(x4, y4).Sub(sp)
	)

	path := func(o *op.Ops) clip.PathSpec {
		var path clip.Path
		path.Begin(p.ctx.Ops)
		path.Move(sp)
		path.Cube(cp0, cp1, ep)
		return path.End()
	}

	stk := op.Push(p.ctx.Ops)
	paint.FillShape(
		p.ctx.Ops,
		rgba(p.cfg.color.stroke),
		clip.Stroke{
			Path:  path(p.ctx.Ops),
			Style: p.cfg.stroke,
		}.Op(),
	)
	stk.Pop()

}

func (p *Proc) poly(ps ...f32.Point) {
	if !p.doFill() && !p.doStroke() {
		return
	}

	path := func(o *op.Ops) clip.PathSpec {
		var path clip.Path
		path.Begin(o)
		path.Move(ps[0])
		for _, p := range ps[1:] {
			path.Line(p.Sub(path.Pos()))
		}
		return path.End()
	}

	if p.doFill() {
		stk := op.Push(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.fill),
			clip.Outline{
				Path: path(p.ctx.Ops),
			}.Op(),
		)
		stk.Pop()
	}

	if p.doStroke() {
		stk := op.Push(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.stroke),
			clip.Stroke{
				Path:  path(p.ctx.Ops),
				Style: p.cfg.stroke,
			}.Op(),
		)
		stk.Pop()
	}
}
