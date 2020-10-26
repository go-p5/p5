// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/go-p5/p5"
)

var (
	numBalls = 12
	spring   = 0.05
	gravity  = 0.03
	friction = -0.9
	balls    = make([]Ball, numBalls)

	width  = 640
	height = 360
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(width, height)
	for i := range balls {
		balls[i] = NewBall(
			rand.Float64()*float64(width),
			rand.Float64()*float64(height),
			rand.Float64()*(70-30)+30,
			i, balls,
		)
	}
	//noStroke();
	p5.Fill(color.RGBA{R: 255, A: 204})
}

func draw() {
	p5.Background(color.Black)
	for i := range balls {
		ball := &balls[i]
		ball.collide()
		ball.move()
		ball.display()
	}
}

type Ball struct {
	x, y   float64
	r      float64
	vx, vy float64

	id    int
	balls []Ball
}

func NewBall(x, y, d float64, i int, balls []Ball) Ball {
	return Ball{
		x: x, y: y, r: 0.5 * d,
		id:    i,
		balls: balls,
	}
}

func (ball *Ball) collide() {
	others := ball.balls[ball.id+1:]
	for i := range others {
		o := &others[i]
		dx := o.x - ball.x
		dy := o.y - ball.y
		dist := math.Hypot(dx, dy)
		minDist := o.r + ball.r
		if dist < minDist {
			angle := math.Atan2(dy, dx)
			sin, cos := math.Sincos(angle)
			tgtX := ball.x + cos*minDist
			tgtY := ball.y + sin*minDist
			ax := (tgtX - o.x) * spring
			ay := (tgtY - o.y) * spring
			ball.vx -= ax
			ball.vy -= ay
			o.vx += ax
			o.vy += ay
		}
	}
}

func (ball *Ball) move() {
	ball.vy += gravity
	ball.x += ball.vx
	ball.y += ball.vy
	switch {
	case ball.x+ball.r > float64(width):
		ball.x = float64(width) - ball.r
		ball.vx *= friction

	case ball.x-ball.r < 0:
		ball.x = ball.r
		ball.vx *= friction
	}

	switch {
	case ball.y+ball.r > float64(height):
		ball.y = float64(height) - ball.r
		ball.vy *= friction

	case ball.y-ball.r < 0:
		ball.y = ball.r
		ball.vy *= friction
	}
}

func (ball *Ball) display() {
	p5.Fill(color.RGBA{R: 255, A: 204})
	p5.Ellipse(ball.x, ball.y, 2*ball.r, 2*ball.r)
}
