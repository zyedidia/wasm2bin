#include <stdio.h>
int main() {
    printf("%p\n", fopen("/home/zyedidia/wasm.txt", "rb"));
    printf("%p\n", fopen("hello.c", "rb"));
    return 0;
}
