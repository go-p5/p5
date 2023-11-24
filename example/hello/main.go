// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"math"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"github.com/go-fonts/latin-modern/lmroman12regular"
	"github.com/go-p5/p5"
)

var (
	fonts = gofont.Collection()
)

func main() {
	loadFonts()
	p5.Run(setup, draw)
}

func setup() {
	p5.LoadFonts(fonts)
	p5.Canvas(400, 400)
	p5.Background(color.Gray{Y: 220})
}

func draw() {
	p5.StrokeWidth(2)
	p5.Fill(color.RGBA{R: 255, A: 208})
	p5.Ellipse(50, 50, 80, 80)

	p5.Fill(color.RGBA{B: 255, A: 208})
	p5.Quad(50, 50, 80, 50, 80, 120, 60, 120)

	p5.Fill(color.RGBA{G: 255, A: 208})
	p5.Rect(200, 200, 50, 100)

	p5.Fill(color.RGBA{G: 255, A: 208})
	p5.Triangle(100, 100, 120, 120, 80, 120)

	p5.TextSize(24)
	p5.Text("Hello, World!", 10, 250)

	p5.TextFont(font.Font{Typeface: "Latin-Modern"})
	p5.TextSize(24)
	p5.Text("Hello, World!", 10, 300)

	p5.TextFont(font.Font{})
	p5.TextSize(24)
	p5.Text("Hello, World!", 10, 350)

	p5.Stroke(color.Black)
	p5.StrokeWidth(5)
	p5.Arc(300, 100, 80, 20, 0, 1.5*math.Pi)
}

func loadFonts() {
	face, err := opentype.Parse(lmroman12regular.TTF)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %+v", err))
	}
	fonts = append(fonts, font.FontFace{
		Font: font.Font{
			Typeface: "Latin-Modern",
		},
		Face: face,
	})
}
