WASI_SYSROOT=/usr/share/wasi-sysroot

hello: hello.wasm
	WASMC_PATH=$(PWD)/.. wasmc --wasm2c -n hello -o $@ $<
hello.wasm: hello.c
	clang -target wasm32-unknown-wasi -O2 $< --sysroot=$(WASI_SYSROOT) -o $@
