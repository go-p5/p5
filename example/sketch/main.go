// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sketch

import (
	"image/color"

	"github.com/go-p5/p5"
)

func Setup() {
	p5.Canvas(400, 400)
	p5.Background(color.Gray{Y: 220})
}

func Draw() {
	p5.Fill(color.RGBA{R: 255, A: 208})
	p5.Ellipse(50, 50, 80, 80)
}
