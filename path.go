// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"fmt"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/x/stroke"
	bstroke "github.com/andybalholm/stroke"
)

func (p *Proc) BeginPath() *Path {
	pp := &Path{proc: p}
	return pp
}

type Path struct {
	proc *Proc
	segs []segment
	vtx  int
}

func (p *Path) pt(x, y float64) f32.Point {
	return p.proc.pt(x, y)
}

func (p *Path) inc() { p.vtx++ }

func (p *Path) Vertex(x, y float64) {
	defer p.inc()
	if p.vtx == 0 {
		p.segs = append(p.segs, opMoveTo(p.pt(x, y)))
		return
	}
	p.segs = append(p.segs, opLineTo(p.pt(x, y)))
}

// Cube draws a cubic Bézier curve from the current position
// to the (x3,y3) point, with the (x1,y1) and (x2,y2) control points.
func (p *Path) Cube(x1, y1, x2, y2, x3, y3 float64) {
	defer p.inc()
	p.segs = append(p.segs, opCubeTo(
		p.pt(x1, y1),
		p.pt(x2, y2),
		p.pt(x3, y3),
	))
}

// Quad draws a quadratic Bézier curve from the current position to
// the (x2,y2) point, with the (x1,y1) control point.
func (p *Path) Quad(x1, y1, x2, y2 float64) {
	defer p.inc()
	p.segs = append(p.segs, opQuadTo(
		p.pt(x1, y1),
		p.pt(x2, y2),
	))
}

// Close closes the current path.
func (p *Path) Close() {
	p.segs = append(p.segs, segment{op: segOpClose})
}

func (p *Path) End() {
	if p.proc.doFill() {
		stk := op.TransformOp{}.Push(p.proc.ctx.Ops)
		paint.FillShape(
			p.proc.ctx.Ops,
			rgba(p.proc.stk.cur().fill),
			p.outline(),
		)
		stk.Pop()
	}

	if p.proc.doStroke() {
		stk := op.TransformOp{}.Push(p.proc.ctx.Ops)
		paint.FillShape(
			p.proc.ctx.Ops,
			rgba(p.proc.stk.cur().stroke.color),
			p.stroke(),
		)
		stk.Pop()
	}

	p.proc = nil
}

func (p *Path) outline() clip.Op {
	return segments(p.segs).outline(p.proc.ctx.Ops)
}

func (p *Path) stroke() clip.Op {
	return segments(p.segs).stroke(p.proc.ctx.Ops, p.proc.stk.cur().stroke)
}

type segment struct {
	op   segmentOp
	args [3]f32.Point
}

type segmentOp uint8

const (
	segOpMoveTo segmentOp = iota
	segOpLineTo
	segOpQuadTo
	segOpCubeTo
	segOpArcTo
	segOpClose
)

func opMoveTo(p f32.Point) segment {
	s := segment{
		op: segOpMoveTo,
	}
	s.args[0] = p
	return s
}

func opLineTo(p f32.Point) segment {
	s := segment{
		op: segOpLineTo,
	}
	s.args[0] = p
	return s
}

func opQuadTo(ctrl, end f32.Point) segment {
	s := segment{
		op: segOpQuadTo,
	}
	s.args[0] = ctrl
	s.args[1] = end
	return s
}

func opCubeTo(ctrl0, ctrl1, end f32.Point) segment {
	s := segment{
		op: segOpCubeTo,
	}
	s.args[0] = ctrl0
	s.args[1] = ctrl1
	s.args[2] = end
	return s
}

func opArcTo(f1, f2 f32.Point, angle float32) segment {
	s := segment{
		op: segOpArcTo,
	}
	s.args[0] = f1
	s.args[1] = f2
	s.args[2].X = angle
	return s
}

type segments []segment

func (segs segments) outline(ops *op.Ops) clip.Op {
	var path clip.Path
	path.Begin(ops)
	for _, seg := range segs {
		switch seg.op {
		case segOpMoveTo:
			path.MoveTo(seg.args[0])
		case segOpLineTo:
			path.LineTo(seg.args[0])
		case segOpArcTo:
			var (
				f1    = seg.args[0]
				f2    = seg.args[1]
				angle = seg.args[2].X
			)
			path.ArcTo(f1, f2, angle)
		case segOpQuadTo:
			var (
				ctl = seg.args[0]
				end = seg.args[1]
			)
			path.QuadTo(ctl, end)
		case segOpCubeTo:
			var (
				ctl0 = seg.args[0]
				ctl1 = seg.args[1]
				end  = seg.args[2]
			)
			path.CubeTo(ctl0, ctl1, end)

		case segOpClose:
			path.Close()

		default:
			panic(fmt.Errorf("p5: unknown outline-path component %d", seg.op))
		}
	}
	return clip.Outline{
		Path: path.End(),
	}.Op()
}

func (segs segments) stroke(ops *op.Ops, sty strokeStyle) clip.Op {
	var (
		shape = stroke.Stroke{
			Width:  sty.style.width,
			Cap:    sty.style.cap,
			Join:   sty.style.join,
			Dashes: sty.style.dashes,
		}
		add = func(seg stroke.Segment) {
			shape.Path.Segments = append(shape.Path.Segments, seg)
		}
		pen f32.Point
		beg f32.Point
	)

	for i, seg := range segs {
		if i == 0 {
			beg = seg.args[0]
		}
		switch seg.op {
		case segOpMoveTo:
			add(stroke.MoveTo(seg.args[0]))
			pen = seg.args[0]
		case segOpLineTo:
			add(stroke.LineTo(seg.args[0]))
			pen = seg.args[0]
		case segOpArcTo:
			arcs := arcTo(pen, seg.args[0], seg.args[1], seg.args[2].X)
			shape.Path.Segments = append(shape.Path.Segments, xStroke(arcs)...)
			pen = f32.Point(arcs[len(arcs)-1].End)
		case segOpQuadTo:
			add(stroke.QuadTo(seg.args[0], seg.args[1]))
			pen = seg.args[1]
		case segOpCubeTo:
			add(stroke.CubeTo(seg.args[0], seg.args[1], seg.args[2]))
			pen = seg.args[2]
		case segOpClose:
			add(stroke.LineTo(beg))
			pen = beg
		default:
			panic(fmt.Errorf("p5: unknown stroke-path component %d", seg.op))
		}
	}

	return shape.Op(ops)
}

func arcTo(start, f1, f2 f32.Point, angle float32) []bstroke.Segment {
	if f1 == f2 {
		return bstroke.AppendArc(nil, bstroke.Pt(start.X, start.Y), bstroke.Pt(f1.X, f1.Y), angle)
	}
	return bstroke.AppendEllipticalArc(nil, bstroke.Pt(start.X, start.Y), bstroke.Pt(f1.X, f1.Y), bstroke.Pt(f2.X, f2.Y), angle)
}

func xStroke(bs []bstroke.Segment) []stroke.Segment {
	vs := make([]stroke.Segment, len(bs))
	for i, b := range bs {
		vs[i] = stroke.CubeTo(f32.Point(b.CP1), f32.Point(b.CP2), f32.Point(b.End))
	}
	return vs
}
