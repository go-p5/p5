// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate embedmd -w README.md

// Package p5 provides types and functions to draw geometrical shapes on
// a canvas, interact with the mouse and keyboard, with the aim to learn
// programming with a graphics aid.
//
// p5 is inspired by Processing and its p5js port:
//   - https://processing.org
//   - https://p5js.org
//
// A very simple p5 program could look like:
//
//	func main() {
//	    p5.Run(setup, draw)
//	}
//
//	func setup() {
//	    p5.Canvas(200, 200)
//	    p5.Background(color.Black)
//	}
//
//	func draw() {
//	    p5.Fill(color.White)
//	    p5.Square(10, 10, 50, 50)
//	}
//
// p5 actually provides two set of APIs:
//   - one closely following the p5js API, with global functions and hidden state,
//   - another one based on the p5.Proc type that encapsulates state.
package p5 // import "github.com/go-p5/p5"

var (
	// gproc is the global Proc instance used by the p5js-like API.
	gproc = newProc(defaultWidth, defaultHeight)
)

// Run executes the user functions setup and draw.
// Run never exits.
func Run(setup, draw Func) {
	gproc.Setup = setup
	gproc.Draw = draw
	gproc.Run()
}

// Func is the type of functions users provide to p5.
type Func func()
