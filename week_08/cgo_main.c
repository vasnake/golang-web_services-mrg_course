#include "_cgo_export.h"

extern void printResultGolang(int i); // declaration of Go function

// implementation of C function
void Multiply(int a, int b) {
    printResultGolang(a*b);
}
