// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
)

// stackOps holds a stack of Gio operations and state.
type stackOps struct {
	ops *op.Ops
	ctx []context
}

func newStackOps(ops *op.Ops) *stackOps {
	return &stackOps{
		ops: ops,
		ctx: make([]context, 1),
	}
}

// context holds the state of the graphics stack.
type context struct {
	bkg    color.Color
	fill   color.Color
	stroke strokeStyle
	text   textStyle

	stk op.StackOp
}

type strokeStyle struct {
	color color.Color
	style clip.StrokeStyle
}

type textStyle struct {
	color color.Color
	align text.Alignment
	size  float32
}

func (stk *stackOps) cur() *context {
	return &stk.ctx[len(stk.ctx)-1]
}

func (stk *stackOps) push() {
	stk.ctx = append(stk.ctx, *stk.cur())
	stk.cur().stk = op.Push(stk.ops)
}

func (stk *stackOps) pop() {
	stk.cur().stk.Pop()
	stk.ctx = stk.ctx[:len(stk.ctx)-1]
}

func (stk *stackOps) rotate(angle float64) {
	ops := stk.ops
	aff := f32.Affine2D{}.Rotate(f32.Pt(0, 0), float32(-angle))
	op.Affine(aff).Add(ops)
}

func (stk *stackOps) scale(x, y float64) {
	ops := stk.ops
	aff := f32.Affine2D{}.Scale(
		f32.Pt(0, 0),
		f32.Pt(float32(x), float32(y)),
	)
	op.Affine(aff).Add(ops)
}

// Push saves the current drawing style settings and transformations.
func (p *Proc) Push() {
	p.stk.push()
}

// Pop restores the previous drawing style settings and transformations.
func (p *Proc) Pop() {
	p.stk.pop()
}

// Rotate rotates the graphical context by angle radians.
// Positive angles rotate counter-clockwise.
func (p *Proc) Rotate(angle float64) {
	p.stk.rotate(angle)
}

// Scale rescales the graphical context by x and y.
func (p *Proc) Scale(x, y float64) {
	p.stk.scale(x, y)
}
