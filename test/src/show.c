#include <stdio.h>
#include <stdlib.h>
int showInt(int x) { printf("%d", x); }
int malloc4(int *x){
    int n, i;
    n = 4;
    x = (int *)malloc(n * sizeof(int)); // 確保したメモリの先頭を引数に代入
    for (i=0; i<n; i++) x[i] = i + 1;
}
