// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package p5

import (
	"log"
	"os"
)

func init() {
	// Enable software rendering for all tests.
	// This allows to have stable results wrt CI.
	err := os.Setenv("LIBGL_ALWAYS_SOFTWARE", "1")
	if err != nil {
		log.Panicf("could not enable GL-software-rendering: %+v", err)
	}
}
