// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/go-p5/p5"
	xdraw "golang.org/x/image/draw"
)

func main() {
	p5.Run(setup, draw)
}

const (
	width  = 640
	height = 640

	ratio = 0.005
	xmax  = width
	ymax  = height
)

var img Image

func setup() {
	p5.Canvas(width, height)

	// FIXME(sbinet): use 'embed' when we no longer support Go-1.15.
	resp, err := http.Get("https://github.com/go-p5/p5/raw/main/testdata/gopher.png")
	if err != nil {
		log.Fatalf("could not fetch image: %+v", err)
	}
	defer resp.Body.Close()

	src, err := png.Decode(resp.Body)
	if err != nil {
		log.Fatalf("could not decode PNG image: %+v", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, width/4, height/4))
	xdraw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)

	img.img = dst
}

func draw() {
	p5.Background(color.Gray{Y: 220})
	img.move()
	img.display()
}

type Image struct {
	img  image.Image
	x, y float64
}

func (img *Image) move() {
	img.x += ratio * xmax
	img.y += ratio * ymax

	if img.x > 1.5*xmax {
		img.x = 0
		img.y = 0
	}
}

func (img *Image) display() {
	f := img.x / xmax
	p5.Push()
	p5.Scale(f, f)
	p5.Translate(f*width/2, f*height/2)
	p5.Rotate(f * math.Pi)
	p5.DrawImage(img.img, 0, 0)
	p5.Pop()

	time.Sleep(10 * time.Millisecond)
}
