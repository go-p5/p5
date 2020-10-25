// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5 // import "github.com/go-p5/p5"

import (
	"log"
	"os"

	"gioui.org/app"
)

var gctx = newContext(400, 400)

// Run executes the user functions setup and draw.
// Run never exits.
func Run(setup, draw Func) {
	go func() {
		err := gctx.run(setup, draw)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// Func is the type of functions users provide to p5.
type Func func()
