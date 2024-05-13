#include "_cgo_export.h"

extern void printResultGolang(int i);

void Multiply(int a, int b) {
    printResultGolang(a*b);
}

int MMultiply(int a, int b) {
    return a * b;
}

#include <unistd.h>
void foo() { }

#include <stdio.h>
#include <stdlib.h>

#include <unistd.h>
