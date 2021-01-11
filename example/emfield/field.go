// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"math"

	"github.com/go-p5/p5"
	"gonum.org/v1/gonum/spatial/r1"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/plot/palette/moreland"
)

const (
	emK     = 1 / (4 * math.Pi * 8.854187 * 10e-12) // vacuum permittivity
	eCharge = 1.602176565 * 10e-19                  // electron charge in coulombs

	ngrad = 1000 // gradient colors
)

var (
	pal = moreland.SmoothBlueRed().Palette(ngrad + 1)
)

type Particle struct {
	p r2.Vec  // p is the particle's position.
	q float64 // q is the particle's charge.
	r float64 // r is the particle's radius.
}

type point struct {
	p r2.Vec
	q float64
	c color.Color
}

type path struct {
	pts []point
	q   float64
	usr bool
}

type EMField struct {
	ps    []Particle
	paths []path

	step uint64
	stop bool

	minField float64
	field    struct {
		r1.Interval
		pos color.Color
		neg color.Color
	}
	r float64 // particle radius
}

func newEMField(ps []Particle, nlines int) *EMField {
	em := EMField{
		ps:       make([]Particle, len(ps)),
		minField: 2 * 10e-10,
		r:        4 * 10e-5,
	}
	em.field.Interval = r1.Interval{
		Min: 9.06e-12, // Gradiant min value
		Max: 1.5e-6,   // Gradiant max value
	}
	em.field.pos = color.RGBA{R: 255, A: 255}
	em.field.neg = color.RGBA{B: 255, A: 255}

	copy(em.ps, ps)

	pts := make([]point, 0, len(ps)*nlines)
	for _, p := range ps {
		for j := 0; j < nlines; j++ {
			sin, cos := math.Sincos(2 * math.Pi / float64(nlines) * float64(j))
			sin *= 10e-5
			cos *= 10e-5
			pts = append(pts, point{
				p: r2.Vec{X: p.p.X + cos, Y: p.p.Y + sin},
				q: p.q,
			})
		}
	}

	em.paths = make([]path, len(pts))
	for i := range pts {
		em.paths[i] = path{
			pts: []point{pts[i]},
			q:   pts[i].q,
		}
	}

	return &em
}

func (em *EMField) update() {
	defer func() { em.step++ }()

	for i := range em.paths {
		em.updatePath(i)
	}
}

func (em *EMField) draw() {
	for i := range em.paths {
		p := &em.paths[i]
		switch {
		case p.usr:
			p5.Fill(nil)
			p5.StrokeWidth(3)
		default:
			p5.Fill(nil)
			p5.StrokeWidth(1)
		}

		if len(p.pts) <= 1 {
			continue
		}

		for j := range p.pts {
			if j == 0 {
				continue
			}
			p0 := p.pts[j-1].p
			p1 := p.pts[j].p
			c1 := p.pts[j].c

			p5.Stroke(c1)
			p5.Line(p0.X, p0.Y, p1.X, p1.Y)
		}
	}

	for _, p := range em.ps {
		c := em.field.pos
		if p.q < 0 {
			c = em.field.neg
		}

		p5.Fill(c)
		p5.Ellipse(p.p.X, p.p.Y, em.r, em.r)
	}
}

func (em *EMField) updatePath(i int) {
	var (
		path     = &em.paths[i]
		last     = path.pts[len(path.pts)-1]
		vec, mag = em.emfield(last.p, path.q)
		fmag     = r2.Norm(vec)
	)

	if !em.stop &&
		(last.p.X < xmin || last.p.X > xmax ||
			last.p.Y < ymin || last.p.Y > ymax) {
		em.stop = true
	}

	for _, p := range em.ps {
		if em.stop && !path.usr && fmag < em.minField {
			//em.paths = append(em.paths[:i], em.paths[i+1:]...)
			return
		}

		switch {
		case fmag > 10e-5:
			vec = vec.Scale(10e-5 / r2.Norm(vec))
		case em.stop && !path.usr && math.Hypot(vec.X-p.p.X, vec.Y-p.p.Y) < p.r:
			//em.paths = append(em.paths[:i], em.paths[i+1:]...)
			return
		}
	}

	path.pts = append(path.pts, point{
		p: vec.Scale(float64(em.step)).Add(last.p),
		c: em.color(mag),
	})
}

func (em *EMField) emfield(p r2.Vec, q float64) (vec r2.Vec, mag float64) {
	// See:
	//  https://en.wikipedia.org/wiki/Electrostatics#Electric_field
	for _, pp := range em.ps {
		v := p.Sub(pp.p)
		m := r2.Norm(v)
		f := 1 / (m * m * m) * emK * q * pp.q * eCharge
		vec = vec.Add(v.Scale(f))
		mag += f * 10e-5
	}
	mag *= q
	return vec, mag
}

func (em *EMField) color(p float64) color.Color {
	v := (p - em.field.Min) / (em.field.Max - em.field.Min)
	switch {
	case p > em.field.Max || v > 1:
		v = 1
	case v < 0:
		v = 0
	}

	i := int(v * ngrad)
	c := pal.Colors()[i]
	return color.NRGBAModel.Convert(c)
}
