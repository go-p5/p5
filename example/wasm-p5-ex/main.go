// Copyright ©2021 The go-p5 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run ./gen.go

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

var (
	listen   = flag.String("addr", ":5555", "server address:port")
	wasmName = "main.wasm"
	wasmExec []byte
	wasmMain []byte
)

func main() {
	flag.Parse()

	initWASM()
	loadWASM()

	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/wasm_exec.js", wasmExecHandle)
	http.HandleFunc("/main.wasm", wasmMainHandle)

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, rootPage, wasmName)
}

const rootPage = `
<html>
	<head>
		<meta charset="utf-8"/>
		<script src="wasm_exec.js"></script>
		<script>
			const go = new Go();
WebAssembly.instantiateStreaming(fetch("%s"), go.importObject).then((result) => {
	go.run(result.instance);
});
		</script>
	</head>
	<body></body>
</html>
`

func wasmExecHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	_, _ = w.Write(wasmExec)
}

func wasmMainHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	_, _ = w.Write(wasmMain)
}

func initWASM() {
	fname := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Panicf("could not locate wasm_exec.js: %+v", err)
	}
	wasmExec = raw
}

func loadWASM() {
	raw, err := ioutil.ReadFile(wasmName)
	if err != nil {
		log.Printf("could not find WASM file: %+v", err)
		log.Fatalf("please run 'go generate'")
	}
	wasmMain = raw
}
