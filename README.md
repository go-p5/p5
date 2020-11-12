# p5

[![GitHub release](https://img.shields.io/github/release/go-p5/p5.svg)](https://github.com/go-p5/p5/releases)
[![go.dev reference](https://pkg.go.dev/badge/github.com/go-p5/p5)](https://pkg.go.dev/github.com/go-p5/p5)
[![CI](https://github.com/go-p5/p5/workflows/CI/badge.svg)](https://github.com/go-p5/p5/actions)
[![GoDoc](https://godoc.org/github.com/go-p5/p5?status.svg)](https://godoc.org/github.com/go-p5/p5)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-p5/p5/blob/main/LICENSE)

`p5` is a simple package that provides primitives resembling the ones exposed by the [p5/processing](https://p5js.org) library.


## Example

```go
package main

import (
	"image/color"

	"github.com/go-p5/p5"
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(400, 400)
	p5.Background(color.Gray{Y: 220})
}

func draw() {
	p5.Fill(color.RGBA{R: 255, A: 208})
	p5.Ellipse(50, 50, 80, 80)

	p5.Fill(color.RGBA{B: 255, A: 208})
	p5.Quad(50, 50, 80, 50, 80, 120, 60, 120)

	p5.Fill(color.RGBA{G: 255, A: 208})
	p5.Rect(200, 200, 50, 100)
}
```
