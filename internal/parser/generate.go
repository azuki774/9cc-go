package parser

import "fmt"

var generatingCode []string // 生成するアセンブリコード
var jumpLabel = 0
var NoMain bool = false // no-main = true なら ソースファイル全体をmain関数とする

func genInitCode() {
	generatingCode = append(generatingCode, ".intel_syntax noprefix\n")
	generatingCode = append(generatingCode, ".globl main\n")
	if NoMain {
		generatingCode = append(generatingCode, "main:\n")
		genCodePrologue("main")
	}

}

func genEndCode() {
	// 関数を呼び出す前にスタックの状態を戻す
	generatingCode = append(generatingCode, "mov rsp, rbp\n")
	generatingCode = append(generatingCode, "pop rbp\n")

	// raw に入っている値を return する
	generatingCode = append(generatingCode, "ret\n")
}

func genLocalVar(node *abstSyntaxNode) {
	// スタックの最後尾に変数のあるアドレスを入れる
	if node.nodeKind != ND_LVAR {
		panic(fmt.Errorf("genLocalVar : left value is not variable"))
	}

	// 変数の値が入っているところにポインタを移動
	generatingCode = append(generatingCode, "mov rax, rbp\n")
	offsetCode := fmt.Sprintf("sub rax, %d\n", node.value.(int)) // offset 分だけずらす
	generatingCode = append(generatingCode, offsetCode)

	generatingCode = append(generatingCode, "push rax\n")
}

func genCodePrologue(funcName string) {
	generatingCode = append(generatingCode, "push rbp\n")
	generatingCode = append(generatingCode, "mov rbp, rsp\n") // rbp のアドレス = rsp のアドレス
	generatingCode = append(generatingCode, "sub rsp, 256\n") // ローカル変数用に容量確保 32 * 8
}

func genCode(node *abstSyntaxNode) {
	// num | local var | = は先に処理する
	nowJumpLabel := jumpLabel
	switch node.nodeKind {
	case ND_NIL:
		return
	case ND_NUM:
		newCode := fmt.Sprintf("push %d\n", node.value.(int))
		generatingCode = append(generatingCode, newCode)
		return
	case ND_LVAR: // local var
		genLocalVar(node)                                           // スタックの最後尾に変数のあるアドレスが入る
		generatingCode = append(generatingCode, "pop rax\n")        // 変数のあるアドレスがスタックから消え、raxに入る
		generatingCode = append(generatingCode, "mov rax, [rax]\n") // rax の中身に rax を書きかえ、変数の値になる
		generatingCode = append(generatingCode, "push rax\n")       // 変数の値をスタックに入れる
		return
	case ND_EQ:
		// 左辺をローカル変数として評価する
		genLocalVar(node.leftNode) // スタックに左辺の変数のアドレスを入れるコード
		genCode(node.rightNode)    // スタックに右辺を計算した結果を入れるコード
		generatingCode = append(generatingCode, "pop rdi\n")
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "mov [rax], rdi\n") // 変数の値を直接右辺に書き換える
		generatingCode = append(generatingCode, "push rdi\n")
		return
	case ND_RETURN:
		genCode(node.leftNode) // return する値を評価するコード
		generatingCode = append(generatingCode, "pop rax\n")
		// スタックを関数呼び出し前に戻す
		generatingCode = append(generatingCode, "mov rsp, rbp\n")
		generatingCode = append(generatingCode, "pop rbp\n")

		generatingCode = append(generatingCode, "ret\n")
		return
	case ND_IF:
		// if (A) B
		jumpLabel++
		genCode(node.leftNode) // A
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		genCode(node.rightNode) // B
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_IFELSE:
		// if (A) B else C
		jumpLabel++
		genCode(node.leftNode) // A
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lelse%d\n", nowJumpLabel))
		genCode(node.rightNode.leftNode) // B
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lend%d\n", nowJumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lelse%d:\n", nowJumpLabel))
		genCode(node.rightNode.rightNode) // C
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_WHILE:
		// while (A) B
		jumpLabel++
		generatingCode = append(generatingCode, fmt.Sprintf(".Lbegin%d:\n", nowJumpLabel))
		genCode(node.leftNode) // A
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		genCode(node.rightNode) // B
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lbegin%d\n", nowJumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_FOR:
		// for (A; B; C) D
		jumpLabel++
		A := node.leftNode.leftNode
		B := node.leftNode.rightNode
		C := node.rightNode.leftNode
		D := node.rightNode.rightNode

		genCode(A)
		generatingCode = append(generatingCode, fmt.Sprintf(".Lbegin%d:\n", nowJumpLabel))
		genCode(B)
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		genCode(D)
		genCode(C)
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lbegin%d\n", nowJumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_BLOCK:
		// { stmt* }
		stmtNodeList := node.value.([]*abstSyntaxNode)
		for _, nowNode := range stmtNodeList {
			genCode(nowNode)
		}
		return
	case ND_FUNDEF:
		funcName := node.value.(string)
		generatingCode = append(generatingCode, fmt.Sprintf("%s:\n", funcName))
		genCodePrologue(funcName)
		genCode(node.rightNode)
		return
	case ND_FUNCALL:
		// value に関数名、leftNode に引数のnode, rightNode に 関数のstmt
		funcName := node.value.(string)

		// rsp を 16の倍数にする調整
		jumpLabel++

		generatingCode = append(generatingCode, "mov rax, rsp\n")
		generatingCode = append(generatingCode, "and rax, 15\n")                          // 下4bitのみにマスキング
		generatingCode = append(generatingCode, fmt.Sprintf("jnz .Lcall%d\n", jumpLabel)) // 下4bit != 0
		generatingCode = append(generatingCode, "mov rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("call %s\n", funcName))
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lend%d\n", jumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lcall%d:\n", jumpLabel)) // 16の倍数になっていなくて、8ずらすときはここから
		generatingCode = append(generatingCode, "sub rsp, 8\n")
		generatingCode = append(generatingCode, "mov rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("call %s\n", funcName))
		generatingCode = append(generatingCode, "add rsp, 8\n")
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", jumpLabel))
		generatingCode = append(generatingCode, "push rax\n")
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
		if NoMain {
			generatingCode = append(generatingCode, "pop rax\n")
		}
	}
	genEndCode()
	return generatingCode, nil
}
