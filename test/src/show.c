#include <stdio.h>
#include <stdlib.h>
int showInt(int x) { printf("%d", x); }
int malloc4(int *x){
    int n;
    n = 4;
    x = (int *)malloc(n * sizeof(int)); // 確保したメモリの先頭を引数に代入
    x[0] = 1;
    x[1] = 2;
    x[2] = 4;
    x[3] = 8;
    return *x;
}
