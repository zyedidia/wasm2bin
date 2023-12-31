package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/pflag"
)

//go:embed rt/w2c2-main.c
var w2c2data []byte

//go:embed rt/uvwasi-main.c
var wasm2cdata []byte

var verbose = pflag.BoolP("verbose", "V", false, "verbose output")

func temp(dir, name string) string {
	tmp, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		log.Fatal(err)
	}
	tmp.Close()
	return tmp.Name()
}

func run(c string) {
	cmd := exec.Command("sh", "-c", c)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	if *verbose {
		fmt.Fprintln(os.Stderr, cmd)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	out := pflag.StringP("output", "o", "a.out", "output file")
	name := pflag.StringP("name", "n", "main", "module name")
	cc := pflag.String("cc", "clang", "C compiler")
	flags := pflag.StringP("flags", "f", "-O2", "C compiler flags")
	useWasm2c := pflag.Bool("wasm2c", false, "Use Wasm2c instead of w2c2")

	pflag.Parse()
	args := pflag.Args()
	if len(args) < 1 {
		log.Fatal("no input")
	}
	in := args[0]

	root := os.Getenv("WASM2BIN")
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if *useWasm2c {
		wabt := filepath.Join(root, "wabt")
		dir, err := os.MkdirTemp("", "root")
		if err != nil {
			log.Fatal(err)
		}

		cwasm := temp(dir, fmt.Sprintf("%s.c", *name))
		hwasm := temp(dir, fmt.Sprintf("%s.h", *name))

		t, err := template.New("uvwas-main").Parse(string(wasm2cdata))
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(filepath.Join(dir, "uvwasi-main.c"))
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(f, struct {
			Hwasm string
			Name  string
		}{
			Hwasm: hwasm,
			Name:  *name,
		})
		if err != nil {
			log.Fatal(err)
		}
		f.Close()

		incwasm2c := filepath.Join(wabt, "wasm2c")
		uvwasirt := filepath.Join(root, "rt", "uvwasi-rt.c")
		incuvwasi := filepath.Join(wabt, "third_party", "uvwasi", "include")
		luvwasi := filepath.Join(wabt, "build", "third_party", "uvwasi")
		luv := filepath.Join(wabt, "build", "_deps", "libuv-build")
		wasmrt := filepath.Join(wabt, "wasm2c", "wasm-rt-impl.c")
		incrt := filepath.Join(root, "rt")
		incsimd := filepath.Join(wabt, "third_party", "simde")

		run(fmt.Sprintf("%s %s -o %s", filepath.Join(wabt, "build", "wasm2c"), in, cwasm))
		run(fmt.Sprintf("%s %s -o %s -I%s -I%s %s -I%s -L%s -luvwasi_a -L%s -luv_a -lm %s %s -I%s %s", *cc, cwasm, *out, incsimd, incwasm2c, uvwasirt, incuvwasi, luvwasi, luv, wasmrt, f.Name(), incrt, *flags))

		os.RemoveAll(dir)
	} else {
		w2c2 := filepath.Join(root, "w2c2")
		dir, err := os.MkdirTemp("", "root")
		if err != nil {
			log.Fatal(err)
		}

		cwasm := temp(dir, fmt.Sprintf("%s.c", *name))
		hwasm := temp(dir, fmt.Sprintf("%s.h", *name))
		incw2c2 := fmt.Sprintf("%s/w2c2", w2c2)
		incwasi := fmt.Sprintf("%s/wasi", w2c2)
		lwasi := fmt.Sprintf("%s/wasi/build", w2c2)
		t, err := template.New("w2c2-main").Parse(string(w2c2data))
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(filepath.Join(dir, "w2c2-main.c"))
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(f, struct {
			Hwasm string
			Name  string
		}{
			Hwasm: hwasm,
			Name:  *name,
		})
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		run(fmt.Sprintf("%s %s %s %s", filepath.Join(w2c2, "build", "w2c2", "w2c2"), in, cwasm, hwasm))
		run(fmt.Sprintf("%s %s -c -o %s.o -I%s -I%s -I%s %s", *cc, cwasm, cwasm, incw2c2, incwasi, pwd, *flags))
		run(fmt.Sprintf("%s %s -c -o %s.o -I%s -I%s -I%s %s", *cc, f.Name(), f.Name(), incw2c2, incwasi, pwd, *flags))
		run(fmt.Sprintf("%s %s.o %s.o -L%s -lw2c2wasi -o %s -lm -static %s", *cc, cwasm, f.Name(), lwasi, *out, *flags))

		os.RemoveAll(dir)
	}
}
