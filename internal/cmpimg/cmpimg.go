// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmpimg compares the raw representation of images taking into account
// idiosyncracies related to their underlying format (SVG, PDF, PNG, ...).
package cmpimg // import "github.com/go-p5/p5/internal/cmpimg"

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"

	_ "image/jpeg"
	_ "image/png"
)

// Equal takes the raw representation of two images, raw1 and raw2,
// together with the underlying image type ("eps", "jpeg", "jpg", "pdf", "png", "svg", "tiff"),
// and returns whether the two images are equal or not.
//
// Equal may return an error if the decoding of the raw image somehow failed.
func Equal(typ string, raw1, raw2 []byte) (bool, error) {
	return EqualApprox(typ, raw1, raw2, 0)
}

// EqualApprox takes the raw representation of two images, raw1 and raw2,
// together with the underlying image type ("eps", "jpeg", "jpg", "pdf", "png", "svg", "tiff"),
// a normalized delta parameter to describe how close the matching should be
// performed (delta=0: perfect match, delta=1, loose match)
// and returns whether the two images are equal or not.
//
// EqualApprox may return an error if the decoding of the raw image somehow failed.
func EqualApprox(typ string, raw1, raw2 []byte, delta float64) (bool, error) {
	switch {
	case delta < 0:
		delta = 0
	case delta > 1:
		delta = 1
	}

	switch typ {
	case "jpeg", "jpg", "png", "tiff":
		v1, _, err := image.Decode(bytes.NewReader(raw1))
		if err != nil {
			return false, err
		}
		v2, _, err := image.Decode(bytes.NewReader(raw2))
		if err != nil {
			return false, err
		}
		if delta == 0 {
			return reflect.DeepEqual(v1, v2), nil
		}
		return cmpImg(v1, v2, delta), nil

	default:
		return false, fmt.Errorf("cmpimg: unknown image type %q", typ)
	}
}

func cmpImg(v1, v2 image.Image, delta float64) bool {
	img1, ok := v1.(*image.RGBA)
	if !ok {
		img1 = newRGBAFrom(v1)
	}

	img2, ok := v2.(*image.RGBA)
	if !ok {
		img2 = newRGBAFrom(v2)
	}

	if len(img1.Pix) != len(img2.Pix) {
		return false
	}

	max := delta * delta
	bnd := img1.Bounds()
	for x := bnd.Min.X; x < bnd.Max.X; x++ {
		for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
			c1 := img1.RGBAAt(x, y)
			c2 := img2.RGBAAt(x, y)
			if !yiqEqApprox(c1, c2, max) {
				return false
			}
		}
	}

	return true
}

// yiqEqApprox compares the colors of 2 pixels, in the NTSC YIQ color space,
// as described in:
//
//	Measuring perceived color difference using YIQ NTSC
//	transmission color space in mobile applications.
//	Yuriy Kotsarenko, Fernando Ramos.
//
// An electronic version is available at:
//
// - http://www.progmat.uaem.mx:8080/artVol2Num2/Articulo3Vol2Num2.pdf
func yiqEqApprox(c1, c2 color.RGBA, d2 float64) bool {
	const max = 35215.0 // difference between 2 maximally different pixels.

	var (
		r1 = float64(c1.R)
		g1 = float64(c1.G)
		b1 = float64(c1.B)

		r2 = float64(c2.R)
		g2 = float64(c2.G)
		b2 = float64(c2.B)

		y1 = r1*0.29889531 + g1*0.58662247 + b1*0.11448223
		i1 = r1*0.59597799 - g1*0.27417610 - b1*0.32180189
		q1 = r1*0.21147017 - g1*0.52261711 + b1*0.31114694

		y2 = r2*0.29889531 + g2*0.58662247 + b2*0.11448223
		i2 = r2*0.59597799 - g2*0.27417610 - b2*0.32180189
		q2 = r2*0.21147017 - g2*0.52261711 + b2*0.31114694

		y = y1 - y2
		i = i1 - i2
		q = q1 - q2

		diff = 0.5053*y*y + 0.299*i*i + 0.1957*q*q
	)
	return diff <= max*d2
}

func newRGBAFrom(src image.Image) *image.RGBA {
	var (
		bnds = src.Bounds()
		dst  = image.NewRGBA(bnds)
	)
	draw.Draw(dst, bnds, src, image.Point{}, draw.Src)
	return dst
}
