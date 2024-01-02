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

	path := func(o *op.Ops, close bool) clip.PathSpec {
		var path clip.Path
		path.Begin(o)
		path.Move(p1)
		path.Arc(f1, f2, 2*math.Pi)
		if close {
			path.Close()
		}
		return path.End()
	}

	if fill := p.stk.cur().fill; fill != nil {
		state := op.Save(p.ctx.Ops)
		close := true
		paint.FillShape(
			p.ctx.Ops,
			rgba(fill),
			clip.Outline{
				Path: path(p.ctx.Ops, close),
			}.Op(),
		)
		state.Load()
	}

	if stroke := p.stk.cur().stroke.color; stroke != nil {
		state := op.Save(p.ctx.Ops)
		close := false
		paint.FillShape(
			p.ctx.Ops,
			rgba(stroke),
			clip.Stroke{
				Path:  path(p.ctx.Ops, close),
				Style: p.stk.cur().stroke.style,
			}.Op(),
		)
		state.Load()
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
		a  = p.cfg.u2sX(w)
		b  = p.cfg.u2sY(h)
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
	defer op.Save(p.ctx.Ops).Load()
	path.Begin(p.ctx.Ops)
	path.Move(p0)
	path.Arc(f1.Sub(p0), f2.Sub(p0), float32(end-beg))

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.stk.cur().stroke.color),
		clip.Stroke{
			Path:  path.End(),
			Style: p.stk.cur().stroke.style,
		}.Op(),
	)
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
	defer op.Save(p.ctx.Ops).Load()
	path.Begin(p.ctx.Ops)
	path.Move(p1)
	path.Line(p2.Sub(path.Pos()))

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.stk.cur().stroke.color),
		clip.Stroke{
			Path:  path.End(),
			Style: p.stk.cur().stroke.style,
		}.Op(),
	)
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

// Bezier draws a cubic Bézier curve from (x1,y1) to (x4,y4) and two control points (x2,y2) and (x3,y3)
func (p *Proc) Bezier(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	if !p.doStroke() {
		return
	}

	var (
		sp   = p.pt(x1, y1)
		cp0  = p.pt(x2, y2).Sub(sp)
		cp1  = p.pt(x3, y3).Sub(sp)
		ep   = p.pt(x4, y4).Sub(sp)
		path clip.Path
	)

	defer op.Save(p.ctx.Ops).Load()

	path.Begin(p.ctx.Ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.stk.cur().stroke.color),
		clip.Stroke{
			Path:  path.End(),
			Style: p.stk.cur().stroke.style,
		}.Op(),
	)
}

// Curve draws a curved line starting at (x2,y2) and ending at (x3,y3).
// (x1,y1) and (x4,y4) are the control points.
//
// Curve is an implementation of Catmull-Rom splines.
func (p *Proc) Curve(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	// Convert the Catmull-Rom curve into a cubic Bézier curve according to
	//
	//  "Conversion Between Bézier and Catmull-Rom Splines"
	//  S. Arasteh and A. Kalisz
	//
	// An electronic version is available at:
	//  https://arxiv.org/pdf/2011.08232.pdf

	tau := p.stk.cur().tau
	if tau == 1 {
		p.Line(x2, y2, x3, y3)
		return
	}

	var (
		cr1 = p.pt(x1, y1)
		cr2 = p.pt(x2, y2)
		cr3 = p.pt(x3, y3)
		cr4 = p.pt(x4, y4)

		itau = 1 / (6 * (1 - tau))

		p1 = cr2
		p2 = cr2.Add(cr3.Sub(cr1).Mul(itau))
		p3 = cr3.Sub(cr4.Sub(cr2).Mul(itau))
		p4 = cr3

		path clip.Path

		sp  = p1
		cp0 = p2.Sub(sp)
		cp1 = p3.Sub(sp)
		ep  = p4.Sub(sp)
	)

	defer op.Save(p.ctx.Ops).Load()

	path.Begin(p.ctx.Ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)

	paint.FillShape(
		p.ctx.Ops,
		rgba(p.stk.cur().stroke.color),
		clip.Stroke{
			Path:  path.End(),
			Style: p.stk.cur().stroke.style,
		}.Op(),
	)
}

// CurveTightness determines how the curve fits to the Curve vertex points.
// CurveTightness controls the Catmull-Rom tau tension.
//
// The default value is 0.
func (p *Proc) CurveTightness(v float64) {
	p.stk.cur().tau = float32(v)
}

func (p *Proc) poly(ps ...f32.Point) {
	if !p.doFill() && !p.doStroke() {
		return
	}

	doClose := false
	if len(ps) > 1 && ps[0] == ps[len(ps)-1] {
		doClose = true
		ps = ps[:len(ps)-1]
	}

	path := func(o *op.Ops) clip.PathSpec {
		var path clip.Path
		path.Begin(o)
		path.Move(ps[0])
		for _, p := range ps[1:] {
			path.Line(p.Sub(path.Pos()))
		}
		if doClose {
			path.Close()
		}
		return path.End()
	}

	if p.doFill() {
		state := op.Save(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.stk.cur().fill),
			clip.Outline{
				Path: path(p.ctx.Ops),
			}.Op(),
		)
		state.Load()
	}

	if p.doStroke() {
		state := op.Save(p.ctx.Ops)
		paint.FillShape(
			p.ctx.Ops,
			rgba(p.stk.cur().stroke.color),
			clip.Stroke{
				Path:  path(p.ctx.Ops),
				Style: p.stk.cur().stroke.style,
			}.Op(),
		)
		state.Load()
	}
}
