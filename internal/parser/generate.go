package parser

import "fmt"

var generatingCode []string // 生成するアセンブリコード

func genInitCode() {
	generatingCode = append(generatingCode, ".intel_syntax noprefix\n")
	generatingCode = append(generatingCode, ".globl main\n")
	generatingCode = append(generatingCode, "main:\n")
}

func genEndCode() {
	generatingCode = append(generatingCode, "pop rax\n")
	generatingCode = append(generatingCode, "ret\n")
}

func genCode(node *abstSyntaxNode) {
	if node.nodeKind == ND_NUM {
		newCode := fmt.Sprintf("push %d\n", node.value.(int))
		generatingCode = append(generatingCode, newCode)
		return
	}

	genCode(node.leftNode)
	genCode(node.rightNode)

	generatingCode = append(generatingCode, "pop rdi\n")
	generatingCode = append(generatingCode, "pop rax\n")

	switch node.nodeKind {
	case ND_ADD:
		generatingCode = append(generatingCode, "add rax, rdi\n")
	case ND_SUB:
		generatingCode = append(generatingCode, "sub rax, rdi\n")
	case ND_MUL:
		generatingCode = append(generatingCode, "imul rax, rdi\n")
	case ND_DIV:
		generatingCode = append(generatingCode, "cqo\n")
		generatingCode = append(generatingCode, "idiv rdi\n")
	}

	generatingCode = append(generatingCode, "push rax\n")
}

func GenAssembleMain(node *abstSyntaxNode) (codes []string, err error) {
	genInitCode()
	genCode(node)
	genEndCode()
	return generatingCode, nil
}
