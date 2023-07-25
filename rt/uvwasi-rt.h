#ifndef UVWASI_RT_H
#define UVWASI_RT_H

#include "uvwasi.h"
#include "wasm-rt.h"

struct w2c_wasi__snapshot__preview1 {
    uvwasi_t * uvwasi;
    wasm_rt_memory_t * instance_memory;
};

#endif
