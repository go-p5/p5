// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5 // import "github.com/go-p5/p5"

var (
	// proc is the global Proc instance used by the p5js-like API.
	proc = newProc(defaultWidth, defaultHeight)
)

// Run executes the user functions setup and draw.
// Run never exits.
func Run(setup, draw Func) {
	proc.Setup = setup
	proc.Draw = draw
	proc.Run()
}

// Func is the type of functions users provide to p5.
type Func func()
