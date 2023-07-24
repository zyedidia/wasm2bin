Wrapper around w2c2 to easily build wasm modules into native binaries.

Example:

```
$ export WASMC_W2C2=/path/to/w2c2
$ wasmc -n hello -o hello hello.wasm
$ ./hello
hello world!
```

# Using wasm2c

Build:

```
cd wasm2c
mkdir build
cd build
cmake .. -G Ninja -DWITH_WASI=ON
ninja wasm2c libuvwasi_a
```
