// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"fmt"
	"image/color"
	"testing"
)

func TestPathVertex(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
			proc.Rect(100, 0, 100, 100)

			proc.StrokeWidth(10)
			proc.Stroke(color.RGBA{B: 255, A: 255})
			p := proc.BeginPath()
			p.Vertex(0, 0)
			p.Vertex(200, 200)
			p.Vertex(0, 200)
			p.Close()
			p.End()
		},
		"testdata/path.png",
	)
	proc.Run(t)
}

func TestPathQuad(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			p := proc.BeginPath()
			p.Vertex(0, 0)
			p.Quad(100, 10, 200, 100)
			p.Quad(100, 10, 100, 200)
			p.Quad(100, 10, 0, 0)
			p.Close()
			p.End()
		},
		"testdata/path_quad.png",
	)
	proc.Run(t)
}

func TestPathCube(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
		},
		func(proc *Proc) {
			p := proc.BeginPath()
			p.Vertex(0, 0)
			p.Cube(50, 10, 150, 10, 200, 100)
			p.Cube(50, 10, 150, 10, 100, 200)
			p.Cube(50, 10, 150, 10, 0, 0)
			p.Close()
			p.End()
		},
		"testdata/path_cube.png",
	)
	proc.Run(t)
}

func TestDraw_Framecount(t *testing.T) {
	const (
		w = 200
		h = 200
	)
	proc := newTestProc(t, w, h,
		func(proc *Proc) {
			proc.Fill(color.RGBA{R: 255, A: 255})
			proc.NoLoop()
		},
		func(proc *Proc) {
			proc.Text(fmt.Sprintf("%d", proc.FrameCount()), 50, 50)
		},
		"testdata/framecount.png",
	)
	proc.Run(t)
}
