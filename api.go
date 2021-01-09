// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import (
	"image/color"
	"log"
)

// Push saves the current drawing style settings and transformations.
func Push() {
	gproc.Push()
}

// Pop restores the previous drawing style settings and transformations.
func Pop() {
	gproc.Pop()
}

// Canvas defines the dimensions of the painting area, in pixels.
func Canvas(w, h int) {
	gproc.Canvas(w, h)
}

// Background defines the background color for the painting area.
// The default color is transparent.
func Background(c color.Color) {
	gproc.Background(c)
}

// Stroke sets the color of the strokes.
func Stroke(c color.Color) {
	gproc.Stroke(c)
}

// StrokeWidth sets the size of the strokes.
func StrokeWidth(v float64) {
	gproc.StrokeWidth(v)
}

// Fill sets the color used to fill shapes.
func Fill(c color.Color) {
	gproc.Fill(c)
}

// TextSize sets the text size.
func TextSize(size float64) {
	gproc.TextSize(size)
}

// Text draws txt on the screen at (x,y).
func Text(txt string, x, y float64) {
	gproc.Text(txt, x, y)
}

// Screenshot saves the current canvas to the provided file.
// Supported file formats are: PNG, JPEG and GIF.
func Screenshot(fname string) {
	err := gproc.Screenshot(fname)
	if err != nil {
		log.Printf("%+v", err)
	}
}

// Rotate rotates the graphical context by angle radians.
// Positive angles rotate counter-clockwise.
func Rotate(angle float64) {
	gproc.Rotate(angle)
}

// Scale rescales the graphical context by x and y.
func Scale(x, y float64) {
	gproc.Scale(x, y)
}
