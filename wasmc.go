package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed rt/rt.c
var rtdata []byte

var verbose = flag.Bool("V", false, "verbose output")

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
	out := flag.String("o", "a.out", "output file")
	name := flag.String("name", "main", "module name")

	flag.Parse()
	args := flag.Args()
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

	cwasm := temp(dir, fmt.Sprintf("%s.c", in))
	owasm := temp(dir, fmt.Sprintf("%s.o", in))
	hwasm := temp(dir, fmt.Sprintf("%s.h", in))
	incw2c2 := fmt.Sprintf("%s/w2c2", w2c2)
	incwasi := fmt.Sprintf("%s/wasi", w2c2)
	lwasi := fmt.Sprintf("%s/wasi/build", w2c2)
	t, err := template.New("rt").Parse(string(rtdata))
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
	run(fmt.Sprintf("clang -O3 %s -c -o %s -I%s", cwasm, owasm, incw2c2))
	run(fmt.Sprintf("clang -O3 %s %s -L%s -lw2c2wasi -I%s -I%s -I%s -o %s", owasm, f.Name(), lwasi, incw2c2, incwasi, pwd, *out))

	os.RemoveAll(dir)
}
