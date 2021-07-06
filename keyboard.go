// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p5

import "gioui.org/io/key"

var KeyIsPressed bool
var Key string

var KeyboardCore struct {
	_DownKeys        map[string]bool
	_KeyReleaseStash map[string]key.Event
	_LastKeyTyped    string

	CallbackKeyPressed  KeyEventFunc
	CallbackKeyTyped    KeyEventFunc
	CallbackKeyReleased KeyEventFunc
}

// KeyIsDown checks if given key is already pressed.
func (p *Proc) KeyIsDown(code string) bool {
	if _, ok := KeyboardCore._DownKeys[code]; ok {
		return true
	}
	return false
}

// KeyEventFunc is the type of key functions users provide to p5.
type KeyEventFunc func(key.Event)

// SetKeyPressedCallback binds the given function to the p5 processor's
// key pressed callback.
func SetKeyPressedCallback(f KeyEventFunc) {
	KeyboardCore.CallbackKeyPressed = f
}

// SetKeyTypedCallback binds the given function to the p5 processor's
// key typed callback.
func SetKeyTypedCallback(f KeyEventFunc) {
	KeyboardCore.CallbackKeyTyped = f
}

// SetKeyReleasedCallback binds the given function to the p5 processor's
// key released callback.
func SetKeyReleasedCallback(f KeyEventFunc) {
	KeyboardCore.CallbackKeyReleased = f
}
