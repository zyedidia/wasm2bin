package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var verbose = flag.Bool("V", false, "verbose output")

func temp(name string) string {
	tmp, err := os.CreateTemp("wasmc", name)
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

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("no input")
	}
	in := args[0]

	w2c2 := os.Getenv("WASMC_W2C2")

	cwasm := temp("*.c")
	owasm := temp("*.o")
	incw2c2 := fmt.Sprintf("-I%s/w2c2", w2c2)
	incwasi := fmt.Sprintf("-I%s/wasi", w2c2)
	lwasi := fmt.Sprintf("-L%s/wasi/build", w2c2)
	run(fmt.Sprintf("w2c2 %s %s", in, cwasm))
	run(fmt.Sprintf("clang -O3 %s -c -o %s -I%s", cwasm, owasm, incw2c2))
	run(fmt.Sprintf("clang -O3 %s %s -L%s -lw2c2wasi -I%s -I%s -o %s", owasm, rt, lwasi, incw2c2, incwasi, *out))
}
