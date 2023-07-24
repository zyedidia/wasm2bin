# Wasmc

Wasmc is a Webassembly compiler that compiles Webassembly modules to native
binaries. It uses either [w2c2](https://github.com/turbolent/w2c2) or
[wasm2c](https://github.com/WebAssembly/wabt/tree/main/wasm2c) to compile
Webassembly to C, and from there links the code with a WASI runtime to produce
a native binary. It uses either w2c2's internal WASI runtime, or
[uvwasi](https://github.com/nodejs/uvwasi).

Example usage:

```
$ export WASMC_PATH=/path/to/wasmc
$ wasmc -n hello -o hello hello.wasm
$ ./hello
hello world!
```

# Installation

First clone the repository and initialize submodules.

```
git clone https://github.com/zyedidia/wasmc
cd wasmc
git submodule update --init --recursive
```

Now choose the backend you would like to use (w2c2 or wasm2c) and build it.

## w2c2

Build w2c2:

```
cd w2c2
cmake -B build
cmake --build build

cmake -B wasi/build
cmake --build wasi/build
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

Note: this uses version 1.30.0 of wasm2c. We do not currently support newer
versions.

# Usage

```
Usage of wasmc:
      --cc string       C compiler (default "clang")
  -f, --flags string    C compiler flags (default "-O2")
  -n, --name string     module name (default "main")
  -o, --output string   output file (default "a.out")
  -V, --verbose         verbose output
      --wasm2c          Use Wasm2c instead of w2c2
```

```
$ export WASMC_PATH=/path/to/wasm
$ wasmc -n hello -o hello hello.wasm          # compile with w2c2
$ wasmc --wasm2c -n hello -o hello hello.wasm # compile with wasm2c
```
