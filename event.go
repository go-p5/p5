// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import "gioui.org/io/key"

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
		Buttons Buttons
	}

	Key struct {
		Cur  key.Event
		Prev key.Event
	}
}

// Buttons is a set of mouse buttons.
type Buttons uint8

// Contain reports whether the set b contains
// all of the buttons.
func (b Buttons) Contain(buttons Buttons) bool {
	return b&buttons == buttons
}

const (
	ButtonLeft Buttons = 1 << iota
	ButtonRight
	ButtonMiddle
)
