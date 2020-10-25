// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

// Event is the current event pushed from the system.
var Event struct {
	Mouse struct {
		Pressed      bool
		PrevPosition struct {
			X float64
			Y float64
		}
		Position struct {
			X float64
			Y float64
		}
	}
}
