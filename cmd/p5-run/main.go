// Copyright ©2020 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "github.com/go-p5/p5/cmd/p5-run"

import (
	"flag"
	"fmt"
	"go/types"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetPrefix("p5-run: ")
	log.SetFlags(0)

	var (
		verbose = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("missing path to sketch")
	}

	pkg, err := importPkg(flag.Arg(0))
	if err != nil {
		log.Fatalf("could not load package %q: %+v", flag.Arg(0), err)
	}

	tmp, err := os.MkdirTemp("", "p5-run-")
	if err != nil {
		log.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	gen := generator{
		pkg:     pkg,
		scope:   "_p5_sketch",
		verbose: *verbose,
		dir:     tmp,
		path:    flag.Arg(0),
	}
	err = gen.generate()
	if err != nil {
		log.Fatalf("could not run generator: %+v", err)
	}

	err = gen.run()
	if err != nil {
		log.Fatalf("could not run p5 sketch: %+v", err)
	}
}

type generator struct {
	pkg   *types.Package
	scope string

	verbose bool
	dir     string // output p5-run directory
	path    string // input import-path or path-to-file
}

func (g *generator) Pkg() string { return g.scope }

func (g *generator) Import() string {
	switch name := g.pkg.Path(); name {
	case "command-line-arguments":
		return fakeModule + "/sketch"
	default:
		return name
	}
}

func (g *generator) Setup() string { return g.get("Setup") }
func (g *generator) Draw() string  { return g.get("Draw") }
func (g *generator) Mouse() string { return g.get("Mouse") }

func (g *generator) get(name string) string {
	obj := g.pkg.Scope().Lookup(name)
	switch obj {
	case nil:
		return "nil"
	default:
		return g.scope + "." + name
	}
}

const (
	fakeModule = "p5-sketch"
	gomod      = "module %s\n"
)

func (g *generator) genGoMod(dir string) error {
	err := os.WriteFile(
		filepath.Join(dir, "go.mod"),
		[]byte(fmt.Sprintf(gomod, fakeModule)), 0644,
	)
	if err != nil {
		return fmt.Errorf("could not generate go.mod: %w", err)
	}

	err = exec.Command("go", "get", "github.com/go-p5/p5").Run()
	if err != nil {
		return fmt.Errorf("could not add p5-require: %w", err)
	}
	return err
}

func (g *generator) generate() error {
	err := g.genGoMod(g.dir)
	if err != nil {
		return fmt.Errorf("could not generate go.mod: %w", err)
	}

	if g.pkg.Path() == "command-line-arguments" {
		err = os.Mkdir(filepath.Join(g.dir, "sketch"), 0755)
		if err != nil {
			return fmt.Errorf("could not create sketch dir: %w", err)
		}
		out, err := os.Create(filepath.Join(g.dir, "sketch", filepath.Base(g.path)))
		if err != nil {
			return fmt.Errorf("could not create sketch file: %w", err)
		}
		defer out.Close()

		f, err := os.Open(g.path)
		if err != nil {
			return fmt.Errorf("could not open input sketch file: %w", err)
		}
		defer f.Close()

		_, err = io.Copy(out, f)
		if err != nil {
			return fmt.Errorf("could not copy input sketch file: %w", err)
		}

		err = out.Close()
		if err != nil {
			return fmt.Errorf("could not save sketch file: %w", err)
		}
	}

	tmpl := template.Must(template.New("sketch").Parse(`package main

import (
	"github.com/go-p5/p5"
	{{.Pkg}} "{{.Import}}"
)

func main() {
	p5.Proc{
		Setup: {{.Setup}},
		Draw:  {{.Draw}},
		Mouse: {{.Mouse}},
	}.Run()
}
`))

	f, err := os.Create(filepath.Join(g.dir, "main.go"))
	if err != nil {
		return fmt.Errorf("could not create sketch code: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, g)
	if err != nil {
		return fmt.Errorf("could not generate sketch code: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not save sketch code: %w", err)
	}

	if g.verbose {
		raw, _ := os.ReadFile(f.Name())
		log.Printf("code:\n%s\n", string(raw))
	}

	args := []string{
		"build", "-o=p5-sketch.exe", "-mod=mod",
	}
	if g.verbose {
		args = append(args, "-v")
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = g.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not build p5-sketch: %+v", err)
	}

	return nil
}

func (g *generator) run() error {
	cmd := exec.Command(filepath.Join(g.dir, "p5-sketch.exe"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run p5-sketch: %+v", err)
	}

	return nil
}

func importPkg(p string) (*types.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, p)
	if err != nil {
		return nil, fmt.Errorf("could not load package %q: %w", p, err)
	}

	return pkgs[0].Types, nil
}
