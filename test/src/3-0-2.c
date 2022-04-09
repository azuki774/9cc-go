int plusfor() {
  int retvalue;
  int i;
  retvalue = 0;
  for (i = 1; i <= 10; i = i + 1){
    retvalue = retvalue + i;
  }
  return retvalue;
}

int main() {
  return plusfor();
}
