WASI_SYSROOT=/usr/share/wasi-sysroot

hello: hello.wasm
	WASM2BIN=$(PWD)/.. wasm2bin --wasm2c --flags "-pthread -ldl" -n hello -o $@ $<
hello.wasm: hello.c
	clang -target wasm32-unknown-wasi -O2 $< --sysroot=$(WASI_SYSROOT) -o $@
