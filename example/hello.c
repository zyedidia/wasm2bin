#include <stdio.h>
int main() {
    printf("hello world\n");
    return 0;
}

void foo(int* p, int i) {
    *p = i;
}
