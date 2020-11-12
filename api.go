// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import "image/color"

// Canvas defines the dimensions of the painting area, in pixels.
func Canvas(w, h int) {
	proc.Canvas(w, h)
}

// Background defines the background color for the painting area.
// The default color is transparent.
func Background(c color.Color) {
	proc.Background(c)
}

// Stroke sets the color of the strokes.
func Stroke(c color.Color) {
	proc.Stroke(c)
}

// Fill sets the color used to fill shapes.
func Fill(c color.Color) {
	proc.Fill(c)
}

// TextSize sets the text size.
func TextSize(size float64) {
	proc.TextSize(size)
}

// Text draws txt on the screen at (x,y).
func Text(txt string, x, y float64) {
	proc.Text(txt, x, y)
}
