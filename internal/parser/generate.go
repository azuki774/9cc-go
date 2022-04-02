package parser

import "fmt"

var generatingCode []string // 生成するアセンブリコード

func genInitCode() {
	generatingCode = append(generatingCode, ".intel_syntax noprefix\n")
	generatingCode = append(generatingCode, ".globl main\n")
	generatingCode = append(generatingCode, "main:\n")
	generatingCode = append(generatingCode, "push rbp\n")
	generatingCode = append(generatingCode, "mov rbp, rsp\n") // rbp のアドレス = rsp のアドレス
	generatingCode = append(generatingCode, "sub rsp, 208\n") // ローカル変数用に容量確保
}

func genEndCode() {
	// 関数を呼び出す前にスタックの状態を戻す
	generatingCode = append(generatingCode, "mov rsp, rbp\n")
	generatingCode = append(generatingCode, "pop rbp\n")

	// raw に入っている値を return する
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
	case ND_COMP:
		generatingCode = append(generatingCode, "cmp rax, rdi\n")
		generatingCode = append(generatingCode, "sete al\n")
		generatingCode = append(generatingCode, "movzb rax, al\n")
	case ND_NOTEQ:
		generatingCode = append(generatingCode, "cmp rax, rdi\n")
		generatingCode = append(generatingCode, "setne al\n")
		generatingCode = append(generatingCode, "movzb rax, al\n")
	case ND_LT:
		generatingCode = append(generatingCode, "cmp rax, rdi\n")
		generatingCode = append(generatingCode, "setl al\n")
		generatingCode = append(generatingCode, "movzb rax, al\n")
	case ND_LTQ:
		generatingCode = append(generatingCode, "cmp rax, rdi\n")
		generatingCode = append(generatingCode, "setle al\n")
		generatingCode = append(generatingCode, "movzb rax, al\n")
	}

	generatingCode = append(generatingCode, "push rax\n")
}

func GenAssembleMain(nodes []*abstSyntaxNode) (codes []string, err error) {
	genInitCode()
	for _, node := range nodes {
		genCode(node)

		// 各式の計算結果をスタックからraxにpop
		generatingCode = append(generatingCode, "pop rax\n")
	}
	genEndCode()
	return generatingCode, nil
}
