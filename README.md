# Wasm2bin

Wasm2bin is a tool that compiles Webassembly modules to native binaries. It
uses either [w2c2](https://github.com/turbolent/w2c2) or
[wasm2c](https://github.com/WebAssembly/wabt/tree/main/wasm2c) to compile
Webassembly to C, and from there calls a C compiler and links the code with a
WASI runtime to produce a native binary. For w2c2 it uses its internal WASI
runtime, and for wasm2c it uses [uvwasi](https://github.com/nodejs/uvwasi).

Example usage:

```
$ export WASM2BIN=/path/to/wasm2bin-src
$ wasm2bin -n hello -o hello hello.wasm
$ ./hello
hello world!
```

# Installation

First clone the repository and initialize submodules.

```
git clone https://github.com/zyedidia/wasm2bin
cd wasm2bin
git submodule update --init --recursive
```

Then build the wrapper tool:

```
go install
```

Now choose the backend you would like to use (w2c2 or wasm2c) and build it.

## w2c2

Build w2c2:

```
cd w2c2
cmake -B build
cmake --build build

cd wasi
cmake -B build
cmake --build build
```

## Wasm2c

Build Wasm2c:

```
cd wasm2c
mkdir build
cd build
cmake .. -G Ninja -DWITH_WASI=ON
ninja wasm2c libuvwasi_a.a
```

Note: uses version 1.0.33 of wasm2c.

# Usage

```
Usage of wasm2bin:
      --cc string       C compiler (default "clang")
  -f, --flags string    C compiler flags (default "-O2")
  -n, --name string     module name (default "main")
  -o, --output string   output file (default "a.out")
  -V, --verbose         verbose output
      --wasm2c          Use Wasm2c instead of w2c2
```

```
$ export WASM2BIN=/path/to/wasm2bin-src
$ wasm2bin -n hello -o hello hello.wasm          # compile with w2c2
$ wasm2bin --wasm2c -n hello -o hello hello.wasm # compile with wasm2c
```

Note:

The wasm2c runtime files have been adapted from the
[wasm2native](https://github.com/vshymanskyy/wasm2native) project and this
[pull request](https://github.com/WebAssembly/wabt/pull/2002).
