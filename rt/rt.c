#include "w2c2_base.h"
#include "wasi.h"
#include "{{ .Hwasm }}"
#include <stdio.h>

extern char** environ;

void
trap(
    Trap trap
) {
    fprintf(stderr, "TRAP: %s\n", trapDescription(trap));
    abort();
}

wasmMemory* wasiMemory(void* instance) {
    return hello_memory((helloInstance*)instance);
}

int
main(int argc, char* argv[]) {
    helloInstance instance;
    helloInstantiate(&instance, NULL);
    if (!wasiInit(argc, argv, environ)) {
        fprintf(stderr, "failed to initialize WASI\n");
        return 1;
    }

    if (!wasiFileDescriptorAdd(-1, "/", NULL)) {
        fprintf(stderr, "failed to add preopen\n");
        return 1;
    }

    hello__start(&instance);
    helloFreeInstance(&instance);

    return 0;
}
