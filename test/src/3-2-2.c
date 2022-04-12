int main(){
    int x;
    int *y;
    y = &x;
    *y = 3;
    showInt(x);
    return 0;
}

