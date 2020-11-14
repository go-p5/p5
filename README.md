# p5

[![GitHub release](https://img.shields.io/github/release/go-p5/p5.svg)](https://github.com/go-p5/p5/releases)
[![go.dev reference](https://pkg.go.dev/badge/github.com/go-p5/p5)](https://pkg.go.dev/github.com/go-p5/p5)
[![CI](https://github.com/go-p5/p5/workflows/CI/badge.svg)](https://github.com/go-p5/p5/actions)
[![GoDoc](https://godoc.org/github.com/go-p5/p5?status.svg)](https://godoc.org/github.com/go-p5/p5)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-p5/p5/blob/main/LICENSE)

`p5` is a simple package that provides primitives resembling the ones exposed by the [p5/processing](https://p5js.org) library.

## License

`p5` is released under the `BSD-3` license.

## Documentation

Documentation for `p5` is served by [GoDev](https://pkg.go.dev/github.com/go-p5/p5).

## Contributing

Guidelines for contributing to [go-p5](https://github.com/go-p5/p5) the same than for the [Go project](https://golang.org/doc/contribute.html#commit_changes).

## Installation

This project relies on [Gio](https://gioui.org) for the graphics parts.
As `Gio` uses system libraries to display graphics, you need to install those for your system/OS for `p5` to work properly.
See [Gio/install](https://gioui.org/doc/install) for details.

## Example

[embedmd]:# (example/hello/main.go go /package main/ $)
```go
package main

import (
	"image/color"
	"math"

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
	p5.Text("Hello, World!", 10, 300)

	p5.Stroke(color.Black)
	p5.StrokeWidth(5)
	p5.Arc(300, 100, 80, 20, 0, 1.5*math.Pi)
}
```

![img-hello](https://github.com/go-p5/p5/raw/main/example/hello/out.png)
