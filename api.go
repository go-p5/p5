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

// PhysCanvas sets the dimensions of the painting area, in pixels, and
// associates physical quantities.
func PhysCanvas(w, h int, xmin, xmax, ymin, ymax float64) {
	gproc.PhysCanvas(w, h, xmin, xmax, ymin, ymax)
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

// Translate applies a translation by x and y.
func Translate(x, y float64) {
	gproc.Translate(x, y)
}

// Shear shears the graphical context by the given x and y angles in radians.
func Shear(x, y float64) {
	gproc.Shear(x, y)
}

// Matrix sets the affine matrix transformation.
func Matrix(a, b, c, d, e, f float64) {
	gproc.Matrix(a, b, c, d, e, f)
}

// RandomSeed changes the sequence of numbers generated by Random.
func RandomSeed(seed uint64) {
	gproc.RandomSeed(seed)
}

// Random returns a pseudo-random number in [min,max).
// Random will produce the same sequence of numbers every time the program runs.
// Use RandomSeed with a seed that changes (like time.Now().UnixNano()) in order to
// produce different sequence of numbers.
func Random(min, max float64) float64 {
	return gproc.Random(min, max)
}

// RandomGaussian returns a random number fitting a Gaussian (normal) distribution.
func RandomGaussian(mean, stdDev float64) float64 {
	return gproc.RandomGaussian(mean, stdDev)
}

// FrameCount returns the number of frames that have been displayed since the program started.
func FrameCount() uint64 {
	return gproc.FrameCount()
}

// By default, p5 continuously executes the code within Draw.
// Loop starts the draw loop again, if it was stopped previously by calling NoLoop.
func Loop() {
	gproc.Loop()
}

// NoLoop stops p5 from continuously executing the code within draw().
func NoLoop() {
	gproc.NoLoop()
}

// IsLooping checks if p5 is continuously executing the code within draw() or not.
func IsLooping() bool {
	return gproc.IsLooping()
}
