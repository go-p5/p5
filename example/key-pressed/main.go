// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/color"

	"gioui.org/io/key"
	"github.com/go-p5/p5"
)

func main() {
	p5.SetKeyPressedCallback(keyPressed)
	p5.SetKeyTypedCallback(keyTyped)
	p5.SetKeyReleasedCallback(keyReleased)
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(400, 400)
	p5.Background(color.Gray{Y: 220})
}

func draw() {
	if p5.KeyIsPressed {
		p5.Stroke(color.Black)
		p5.Fill(color.RGBA{R: 255, A: 255})
	} else {
		p5.Stroke(nil)
		p5.Fill(color.Transparent)
	}
	p5.Ellipse(
		200, 200,
		100, 100,
	)

	p5.TextSize(24)
	p5.Text(fmt.Sprintf("count=%d", cnt), 10, 390)
}

var cnt int

func keyPressed(e key.Event) {
	fmt.Println("Key pressed:", e.Name)
	cnt++
}

func keyTyped(e key.Event) {
	fmt.Println("Key typed:", e.Name)
}

func keyReleased(e key.Event) {
	fmt.Println("Key released:", e.Name)
}
