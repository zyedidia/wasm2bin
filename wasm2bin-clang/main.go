package main

// This is a wrapper around clang that invokes a WebAssembly compiler after
// building so that it produces an executable rather than a Wasm binary.

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func run(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	log.Println(c)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var out string
	clang := "clang"
	link := true
	args := make([]string, 0, len(os.Args))
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasSuffix(arg, ".c") ||
			strings.HasSuffix(arg, ".cc") ||
			strings.HasSuffix(arg, ".cpp") ||
			strings.HasSuffix(arg, ".cxx") ||
			strings.HasSuffix(arg, ".c++") ||
			strings.HasSuffix(arg, ".s") ||
			strings.HasSuffix(arg, ".S") ||
			strings.HasSuffix(arg, ".C") {
			link = false
		}
		switch arg {
		case "-compiler":
			if i+1 >= len(os.Args) {
				log.Fatal("-o needs an argument")
			}
			clang = os.Args[i+1]
			i++
		case "-c":
			link = false
			args = append(args, arg)
		case "-o":
			if i+1 >= len(os.Args) {
				log.Fatal("-o needs an argument")
			}
			out = os.Args[i+1]
			i++
			args = append(args, "-o", out)
		default:
			args = append(args, arg)
		}
	}

	if out == "" && link {
		out = "a.out"
	}

	run(clang, args...)
	if !link {
		return
	}

	postlink := os.Getenv("POSTLINKCMD")
	var flags string
	if strings.Contains(postlink, "wasm2c") {
		module := strings.ReplaceAll(out, "_", "__") + "__base"
		flags = fmt.Sprintf("-n %s", module)
	} else if strings.Contains(postlink, "wamrc") {
		flags = "--target=aarch64"
	} else if strings.Contains(postlink, "w2c2") {
		module := strings.ReplaceAll(out, "_", "") + "base"
		flags = fmt.Sprintf("-n %s", module)
	} else if strings.Contains(postlink, "wasmtime") {
		// nothing to do
	}
	in := out + "_base.wasm"
	opt := out + "_base.opt.wasm"
	run("cp", out, in)

	run("sh", "-c", fmt.Sprintf("wasm-opt -O3 %s -o %s; %s %s -o %s %s", in, opt, postlink, flags, out, in))
}
