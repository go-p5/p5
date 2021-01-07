// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

// Ellipse draws an ellipse at (x,y) with the provided width and height.
func Ellipse(x, y, w, h float64) {
	gproc.Ellipse(x, y, w, h)
}

// Circle draws a circle at (x,y) with a diameter d.
func Circle(x, y, d float64) {
	gproc.Circle(x, y, d)
}

// Arc draws an ellipsoidal arc centered at (x,y), with the provided
// width and height, and a path from the beg to end radians.
// Positive angles denote a counter-clockwise path.
func Arc(x, y, w, h float64, beg, end float64) {
	gproc.Arc(x, y, w, h, beg, end)
}

// Line draws a line between (x1,y1) and (x2,y2).
func Line(x1, y1, x2, y2 float64) {
	gproc.Line(x1, y1, x2, y2)
}

// Quad draws a quadrilateral, connecting the 4 points (x1,y1),
// (x2,y2), (x3,y3) and (x4,y4) together.
func Quad(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	gproc.Quad(x1, y1, x2, y2, x3, y3, x4, y4)
}

// Rect draws a rectangle at (x,y) with width w and height h.
func Rect(x, y, w, h float64) {
	gproc.Rect(x, y, w, h)
}

// Square draws a square at (x,y) with size s.
func Square(x, y, s float64) {
	gproc.Square(x, y, s)
}

// Triangle draws a triangle, connecting the 3 points (x1,y1), (x2,y2)
// and (x3,y3) together.
func Triangle(x1, y1, x2, y2, x3, y3 float64) {
	gproc.Triangle(x1, y1, x2, y2, x3, y3)
}

// Bezier draws cubic Bézier curve from (x1,y1), through two control points (x2,y2) and (x3,y3) to (x4,y5).
func Bezier(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	gproc.Bezier(x1, y1, x2, y2, x3, y3, x4, y4)
}
