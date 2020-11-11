// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"time"

	"github.com/go-p5/p5"
	"gonum.org/v1/gonum/spatial/r2"
)

func main() {
	p5.Run(setup, draw)
}

const (
	width  = 1200
	height = 1200
)

func setup() {
	p5.Canvas(width, height)
	p5.Background(color.Black)
}

var (
	tick  time.Time
	first = true
)

func draw() {
	if first {
		first = false
		tick = time.Now()
		return
	}
	var dt = time.Since(tick).Seconds() * 5 * 10e5
	tick = time.Now()

	for i := range sol {
		p := &sol[i]
		p.update(dt, sol)
		p.draw()
	}
}

var (
	sol = []NBody{
		NewNBody(
			"Mercury", 3.3e23,
			r2.Vec{X: 0, Y: 4.7e10},
			r2.Vec{X: 5.9e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 241, G: 203, B: 131, A: 255},
		),
		NewNBody(
			"Venus", 4.9e24,
			r2.Vec{X: 0, Y: 1.1e11},
			r2.Vec{X: 3.5e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 243, G: 223, B: 107, A: 255},
		),
		NewNBody(
			"Earth", 6.0e24,
			r2.Vec{X: 0, Y: 1.5e11},
			r2.Vec{X: 3.0e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 173, G: 231, B: 247, A: 255},
		),
		NewNBody(
			"Mars", 6.4e23,
			r2.Vec{X: 0, Y: 2.1e11},
			r2.Vec{X: 2.6e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 223, G: 120, B: 036, A: 255},
		),
		NewNBody(
			"Jupiter", 1.9e27,
			r2.Vec{X: 0, Y: 7.4e11},
			r2.Vec{X: 1.3e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 243, G: 131, B: 239, A: 255},
		),
		NewNBody(
			"Saturn", 5.6e26,
			r2.Vec{X: 0, Y: 1.3e12},
			r2.Vec{X: 1.0e4, Y: 0},
			r2.Vec{},
			color.RGBA{R: 118, G: 064, B: 045, A: 255},
		),
		NewNBody(
			"Uranus", 8.7e25,
			r2.Vec{X: 0, Y: 2.7e12},
			r2.Vec{X: 7.1e3, Y: 0},
			r2.Vec{},
			color.RGBA{R: 157, G: 221, B: 250, A: 255},
		),
		NewNBody(
			"Neptun", 1.0e26,
			r2.Vec{X: 0, Y: 4.4e12},
			r2.Vec{X: 5.5e3, Y: 0},
			r2.Vec{},
			color.RGBA{R: 045, G: 86, B: 148, A: 255},
		),
		NewNBody(
			"Sun", 1.989e30,
			r2.Vec{}, r2.Vec{}, r2.Vec{},
			color.RGBA{R: 246, G: 244, B: 129, A: 255},
		),
	}
)

type NBody struct {
	mass float64
	pos  r2.Vec
	vel  r2.Vec
	acc  r2.Vec
	path []r2.Vec

	c    color.Color
	name string
}

const (
	orbitLen = 130

	xmin = -1.3e12
	xmax = +1.3e12
	ymin = -1.3e12
	ymax = +1.3e12
)

func tr(x, xmin, xmax float64) float64 {
	return (x - xmin) / (xmax - xmin) * width
}

func NewNBody(name string, mass float64, pos, vel, acc r2.Vec, c color.Color) NBody {
	return NBody{
		mass: mass,
		pos:  pos,
		vel:  vel,
		acc:  acc,
		path: make([]r2.Vec, 0, orbitLen),
		c:    c,
		name: name,
	}
}

func (p *NBody) update(dt float64, ps []NBody) {
	p.updateAcc(dt, ps)
	p.pos = p.pos.Add(p.vel.Scale(dt)).Add(p.acc.Scale(0.5 * dt * dt))

	acc := p.acc

	p.updateAcc(dt, ps)
	p.vel = p.vel.Add(acc.Add(p.acc).Scale(0.5 * dt))
}

func (p *NBody) updateAcc(dt float64, ps []NBody) {
	const G = 6.67430e-11 // Gravitational constant
	p.acc = r2.Vec{}

	for i := range ps {
		q := &ps[i]
		if q.name == p.name {
			continue
		}

		// acceleration for each body.
		var (
			delta = p.pos.Sub(q.pos)
			d     = r2.Norm(delta)
		)
		p.acc = p.acc.Add(delta).Scale(-G * q.mass / (d * d * d))
	}
}

func (p *NBody) draw() {
	r := 12.0

	if p.name != "Sun" {
		r = 6

		// draw orbits
		p.path = append(p.path, p.pos)
		if len(p.path) >= orbitLen {
			copy(p.path, p.path[1:])
			p.path = p.path[:orbitLen]
		}
		for i := 0; i < len(p.path)-1; i++ {
			pi := p.path[i]
			pj := p.path[i+1]

			pix := tr(pi.X, xmin, xmax)
			piy := tr(pi.Y, ymin, ymax)
			pjx := tr(pj.X, xmin, xmax)
			pjy := tr(pj.Y, ymin, ymax)
			p5.Stroke(p.c)
			p5.Line(pix, piy, pjx, pjy)
		}
	}

	px := tr(p.pos.X, xmin, xmax)
	py := tr(p.pos.Y, ymin, ymax)

	p5.Fill(p.c)
	p5.Ellipse(px, py, r, r)
}
