// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func (p *Proc) BeginPath() *Path {
	pp := &Path{proc: p}
	return pp
}

type Path struct {
	proc  *Proc
	funcs []func(p *clip.Path)
	vtx   int
}

func (p *Path) pt(x, y float64) f32.Point {
	return p.proc.pt(x, y)
}

func (p *Path) inc() { p.vtx++ }

func (p *Path) Vertex(x, y float64) {
	defer p.inc()
	if p.vtx == 0 {
		p.funcs = append(p.funcs, func(path *clip.Path) {
			pt := p.pt(x, y).Sub(path.Pos())
			path.Move(pt)
		})
		return
	}
	p.funcs = append(p.funcs, func(path *clip.Path) {
		pt := p.pt(x, y).Sub(path.Pos())
		path.Line(pt)
	})
}

// Cube draws a cubic Bézier curve from the current position
// to the (x3,y3) point, with the (x1,y1) and (x2,y2) control points.
func (p *Path) Cube(x1, y1, x2, y2, x3, y3 float64) {
	defer p.inc()
	p.funcs = append(p.funcs, func(path *clip.Path) {
		var (
			pos  = path.Pos()
			ctl1 = p.pt(x1, y1).Sub(pos)
			ctl2 = p.pt(x2, y2).Sub(pos)
			end  = p.pt(x3, y3).Sub(pos)
		)
		path.Cube(ctl1, ctl2, end)
	})
}

// Quad draws a quadratic Bézier curve from the current position to
// the (x2,y2) point, with the (x1,y1) control point.
func (p *Path) Quad(x1, y1, x2, y2 float64) {
	defer p.inc()
	p.funcs = append(p.funcs, func(path *clip.Path) {
		var (
			pos = path.Pos()
			ctl = p.pt(x1, y1).Sub(pos)
			end = p.pt(x2, y2).Sub(pos)
		)
		path.Quad(ctl, end)
	})
}

// Close closes the current path.
func (p *Path) Close() {
	p.funcs = append(p.funcs, func(path *clip.Path) {
		path.Close()
	})
}

func (p *Path) End() {
	if p.proc.doFill() {
		stk := op.Push(p.proc.ctx.Ops)
		paint.FillShape(
			p.proc.ctx.Ops,
			rgba(p.proc.cfg.color.fill),
			clip.Outline{
				Path: p.path(),
			}.Op(),
		)
		stk.Pop()
	}

	if p.proc.doStroke() {
		stk := op.Push(p.proc.ctx.Ops)
		paint.FillShape(
			p.proc.ctx.Ops,
			rgba(p.proc.cfg.color.stroke),
			clip.Stroke{
				Path:  p.path(),
				Style: p.proc.cfg.stroke,
			}.Op(),
		)
		stk.Pop()
	}

	p.proc = nil
}

func (p *Path) path() clip.PathSpec {
	var path clip.Path
	path.Begin(p.proc.ctx.Ops)
	for _, fct := range p.funcs {
		fct(&path)
	}
	return path.End()
}
