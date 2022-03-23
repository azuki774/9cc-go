.intel_syntax noprefix
.globl main
main:
push 3
push 6
pop rdi
pop rax
add rax, rdi
push rax
pop rax
ret
