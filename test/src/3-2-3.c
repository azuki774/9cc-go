int pnumAdd(int *x){
    *x = *x + 1;
    return 0;
}

int main(){
    int x;
    int *y;
    x = 5;
    y = &x;
    pnumAdd(y);
    showInt(x);
    return 0;
}

