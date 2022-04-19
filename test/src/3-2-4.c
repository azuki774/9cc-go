int main(){
    int **pp;
    int *p;
    int a;
    a = 100;

    pp = &p;
    p = &a;
    showInt(**pp);
}

