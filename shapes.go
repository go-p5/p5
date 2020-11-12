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

	path := func() *clip.Path {
		var path clip.Path
		path.Begin(p.ctx.Ops)
		path.Move(p1)
		path.Arc(f1, f2, 2*math.Pi)
		return &path
	}

	if p.cfg.color.fill != nil {
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.fill),
			path().Outline(),
		)
	}

	if p.cfg.color.stroke != nil {
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.stroke),
			path().Stroke(p.cfg.linew, p.cfg.sty),
		)
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
		c   = p.pt(x, y)
		a   = p.cfg.trX(w)
		b   = p.cfg.trY(h)
		foc = 0.0
		f1  f32.Point
		f2  f32.Point
	)

	switch {
	case a >= b:
		foc = math.Sqrt(a*a - b*b)
		f1 = c.Add(p.pt(+foc, 0))
		f2 = c.Add(p.pt(-foc, 0))
	default:
		foc = math.Sqrt(b*b - a*a)
		f1 = c.Add(p.pt(0, +foc))
		f2 = c.Add(p.pt(0, -foc))
	}

	var (
		sin, cos = math.Sincos(beg)
		p0       = p.pt(a*cos, b*sin).Add(c)
		path     clip.Path
	)
	path.Begin(p.ctx.Ops)
	path.Move(p0)
	path.Arc(f1.Sub(p0), f2.Sub(p0), float32(end-beg))

	arc := path.Stroke(p.cfg.linew, p.cfg.sty)
	paint.FillShape(p.ctx.Ops, rgba(p.cfg.color.stroke), arc)
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
	path.Begin(p.ctx.Ops)
	path.Move(p1)
	path.Line(p2.Sub(path.Pos()))

	line := path.Stroke(p.cfg.linew, p.cfg.sty)

	paint.FillShape(p.ctx.Ops, rgba(p.cfg.color.stroke), line)
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
		proc.pt(x1, y1),
		proc.pt(x2, y2),
		proc.pt(x3, y3),
		proc.pt(x1, y1),
	)
}

func (p *Proc) poly(ps ...f32.Point) {
	if !p.doFill() && !p.doStroke() {
		return
	}

	path := func() *clip.Path {
		var path clip.Path
		path.Begin(p.ctx.Ops)
		path.Move(ps[0])
		for _, p := range ps[1:] {
			path.Line(p.Sub(path.Pos()))
		}
		return &path
	}

	if p.doFill() {
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.fill),
			path().Outline(),
		)
	}

	if p.doStroke() {
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.cfg.color.stroke),
			path().Stroke(p.cfg.linew, p.cfg.sty),
		)
	}
}
