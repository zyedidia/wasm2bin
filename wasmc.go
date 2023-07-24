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

	pflag.Parse()
	args := pflag.Args()
	if len(args) < 1 {
		log.Fatal("no input")
	}
	in := args[0]

	w2c2 := os.Getenv("WASMC_W2C2")
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir := os.TempDir()
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cwasm := temp(dir, fmt.Sprintf("%s.c", *name))
	owasm := temp(dir, fmt.Sprintf("%s.o", *name))
	hwasm := temp(dir, fmt.Sprintf("%s.h", *name))
	incw2c2 := fmt.Sprintf("%s/w2c2", w2c2)
	incwasi := fmt.Sprintf("%s/wasi", w2c2)
	lwasi := fmt.Sprintf("%s/wasi/build", w2c2)
	t, err := template.New("rt").Parse(string(w2c2data))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(filepath.Join(dir, "rt.c"))
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
	run(fmt.Sprintf("w2c2 %s %s %s", in, cwasm, hwasm))
	run(fmt.Sprintf("%s -fomit-frame-pointer -O3 %s -c -o %s -I%s -lm -static", *cc, cwasm, owasm, incw2c2))
	run(fmt.Sprintf("%s -fomit-frame-pointer -O3 %s %s -L%s -lw2c2wasi -I%s -I%s -I%s -o %s -lm -static", *cc, owasm, f.Name(), lwasi, incw2c2, incwasi, pwd, *out))

	os.RemoveAll(dir)
}
