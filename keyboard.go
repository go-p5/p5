// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import "gioui.org/io/key"

var Keyboard struct {
	KeyIsPressed bool
	Key          string

	downKeys        map[string]struct{}
	keyReleaseStash map[string]key.Event
	lastKeyTyped    string

	keyPressedCb  KeyEventFunc
	keyTypedCb    KeyEventFunc
	keyReleasedCb KeyEventFunc
}

// KeyIsDown checks if given key is already pressed.
func (p *Proc) KeyIsDown(code string) bool {
	if _, ok := Keyboard.downKeys[code]; ok {
		return true
	}
	return false
}

// KeyEventFunc is the type of key functions users provide to p5.
type KeyEventFunc func(key.Event)

// SetKeyPressedCallback binds the given function to the p5 processor's
// key pressed callback.
func SetKeyPressedCallback(f KeyEventFunc) {
	Keyboard.keyPressedCb = f
}

// SetKeyTypedCallback binds the given function to the p5 processor's
// key typed callback.
func SetKeyTypedCallback(f KeyEventFunc) {
	Keyboard.keyTypedCb = f
}

// SetKeyReleasedCallback binds the given function to the p5 processor's
// key released callback.
func SetKeyReleasedCallback(f KeyEventFunc) {
	Keyboard.keyReleasedCb = f
}
