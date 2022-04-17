package parser

import "fmt"

var generatingCode []string // 生成するアセンブリコード
var jumpLabel = 0
var NoMain bool = false // no-main = true なら ソースファイル全体をmain関数とする
var argsRegisterName = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func genInitCode() {
	generatingCode = append(generatingCode, ".intel_syntax noprefix\n")
	generatingCode = append(generatingCode, ".globl main\n")
	if NoMain {
		generatingCode = append(generatingCode, "main:\n")
		genCodePrologue("main", 0)
	}

}

func genEndCode() {
	// 関数を呼び出す前にスタックの状態を戻す
	generatingCode = append(generatingCode, "mov rsp, rbp\n")
	generatingCode = append(generatingCode, "pop rbp\n")

	// raw に入っている値を return する
	generatingCode = append(generatingCode, "ret\n")
}

func genLocalVar(node *abstSyntaxNode) (err error) {
	// スタックの最後尾に変数のあるアドレスを入れる

	if node.nodeKind == ND_LVAR {
		// 変数の値が入っているところにポインタを移動
		generatingCode = append(generatingCode, "mov rax, rbp\n")
		offsetCode := fmt.Sprintf("sub rax, %d\n", node.value.(variable).offset) // offset 分だけずらす
		generatingCode = append(generatingCode, offsetCode)

		generatingCode = append(generatingCode, "push rax\n")
		return
	}

	if node.nodeKind == ND_DEREF {
		// 左辺値がDEREF(*)のとき
		genCode(node.leftNode)
		return
	}

	panic(fmt.Errorf("genLocalVar : left value is not variable, actual = %s", node.nodeKind))
}

func genCodePrologue(funcName string, argsNum int) {
	generatingCode = append(generatingCode, "push rbp\n")
	generatingCode = append(generatingCode, "mov rbp, rsp\n") // rbp のアドレス = rsp のアドレス
	generatingCode = append(generatingCode, "sub rsp, 256\n") // ローカル変数用に容量確保 32 * 8
	if funcName == "main" {
		return
	}

	if argsNum > 0 {
		generatingCode = append(generatingCode, "mov rsp, rbp\n") // 一旦変数領域のスタックのベースに移動
		for i := 0; i < argsNum; i++ {
			// 受け取った引数に変数を書き換えていく
			generatingCode = append(generatingCode, "sub rsp, 8\n")                                      // i個目の引数のアドレスに移動
			generatingCode = append(generatingCode, fmt.Sprintf("mov [rsp], %s\n", argsRegisterName[i])) // i個目の引数を書き換え
		}
		generatingCode = append(generatingCode, "mov rsp, rbp\n")
		generatingCode = append(generatingCode, "sub rsp, 256\n") // 元の位置に復帰
	}
}

func genCode(node *abstSyntaxNode) (err error) {
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
		if err := genLocalVar(node); err != nil { // スタックの最後尾に変数のあるアドレスが入る
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")        // 変数のあるアドレスがスタックから消え、raxに入る
		generatingCode = append(generatingCode, "mov rax, [rax]\n") // rax の中身に rax を書きかえ、変数の値になる
		generatingCode = append(generatingCode, "push rax\n")       // 変数の値をスタックに入れる
		return
	case ND_EQ:
		// 左辺をローカル変数として評価する
		if err := genLocalVar(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}
		if err := genCode(node.rightNode); err != nil { // スタックに右辺を計算した結果を入れるコード
			return err
		}
		generatingCode = append(generatingCode, "pop rdi\n")
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "mov [rax], rdi\n") // 変数の値を直接右辺に書き換える
		generatingCode = append(generatingCode, "push rdi\n")
		return
	case ND_ADDR: // &x
		if err := genLocalVar(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}
		return
	case ND_DEREF: // *x
		if err := genCode(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "mov rax, [rax]\n")
		generatingCode = append(generatingCode, "push rax\n")
		return
	case ND_RETURN:
		if err := genCode(node.leftNode); err != nil { // return する値を評価するコード
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		// スタックを関数呼び出し前に戻す
		generatingCode = append(generatingCode, "mov rsp, rbp\n")
		generatingCode = append(generatingCode, "pop rbp\n")

		generatingCode = append(generatingCode, "ret\n")
		return
	case ND_IF:
		// if (A) B
		jumpLabel++
		if err := genCode(node.leftNode); err != nil { // A
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		if err := genCode(node.rightNode); err != nil { // B
			return err
		}
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_IFELSE:
		// if (A) B else C
		jumpLabel++
		if err := genCode(node.leftNode); err != nil { // A
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lelse%d\n", nowJumpLabel))
		if err := genCode(node.rightNode.leftNode); err != nil { // B
			return err
		}
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lend%d\n", nowJumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lelse%d:\n", nowJumpLabel))
		if err := genCode(node.rightNode.rightNode); err != nil { // C
			return err
		}
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_WHILE:
		// while (A) B
		jumpLabel++
		generatingCode = append(generatingCode, fmt.Sprintf(".Lbegin%d:\n", nowJumpLabel))
		if err := genCode(node.leftNode); err != nil { // A
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		if err := genCode(node.rightNode); err != nil { // B
			return err
		}
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

		if err := genCode(A); err != nil {
			return err
		}
		generatingCode = append(generatingCode, fmt.Sprintf(".Lbegin%d:\n", nowJumpLabel))
		if err := genCode(B); err != nil {
			return err
		}
		generatingCode = append(generatingCode, "pop rax\n")
		generatingCode = append(generatingCode, "cmp rax, 0\n")
		generatingCode = append(generatingCode, fmt.Sprintf("je  .Lend%d\n", nowJumpLabel))
		if err := genCode(D); err != nil {
			return err
		}
		if err := genCode(C); err != nil {
			return err
		}
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lbegin%d\n", nowJumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", nowJumpLabel))
		return
	case ND_BLOCK:
		// { stmt* }
		stmtNodeList := node.value.([]*abstSyntaxNode)
		for _, nowNode := range stmtNodeList {
			if err := genCode(nowNode); err != nil {
				return err
			}
		}
		return
	case ND_FUNDEF:
		funcName := node.value.(string)
		generatingCode = append(generatingCode, fmt.Sprintf("%s:\n", funcName))

		// この変数が何変数関数か調べる
		var argsNum int = 0
		argsNode := node.leftNode
		if argsNode.value != nil {
			argsNum = len(argsNode.value.([]*abstSyntaxNode))
		}

		genCodePrologue(funcName, argsNum)
		if err := genCode(node.rightNode); err != nil {
			return err
		}
		return
	case ND_FUNDEF_ARGS:
		argsNodes := node.value.([]*abstSyntaxNode)
		for _, v := range argsNodes {
			if err := genCode(v); err != nil { // 変数を定義する
				return err
			}
		}
		return
	case ND_FUNCALL:
		// value に関数名、leftNode に引数のnode, rightNode に 関数のstmt
		funcName := node.value.(string)

		// ND_FUNCALL_ARGS
		argsNode := node.leftNode.value.([]*abstSyntaxNode)
		for i, v := range argsNode {
			if err := genCode(v); err != nil {
				return err
			}
			generatingCode = append(generatingCode, fmt.Sprintf("pop %s\n", argsRegisterName[i]))
		}
		// rsp を 16の倍数にする調整
		jumpLabel++

		generatingCode = append(generatingCode, "mov rax, rsp\n")
		generatingCode = append(generatingCode, "and rax, 15\n")                          // 下4bitのみにマスキング
		generatingCode = append(generatingCode, fmt.Sprintf("jnz .Lcall%d\n", jumpLabel)) // 下4bit != 0
		generatingCode = append(generatingCode, "mov rax, 0\n")                           // TODO: rax は引数の数
		generatingCode = append(generatingCode, fmt.Sprintf("call %s\n", funcName))
		generatingCode = append(generatingCode, fmt.Sprintf("jmp .Lend%d\n", jumpLabel))
		generatingCode = append(generatingCode, fmt.Sprintf(".Lcall%d:\n", jumpLabel)) // 16の倍数になっていなくて、8ずらすときはここから
		generatingCode = append(generatingCode, "sub rsp, 8\n")
		generatingCode = append(generatingCode, "mov rax, 0\n") // TODO: rax は引数の数
		generatingCode = append(generatingCode, fmt.Sprintf("call %s\n", funcName))
		generatingCode = append(generatingCode, "add rsp, 8\n")
		generatingCode = append(generatingCode, fmt.Sprintf(".Lend%d:\n", jumpLabel))
		generatingCode = append(generatingCode, "push rax\n")
		return
	}

	// 二項演算系

	if err := genCode(node.leftNode); err != nil {
		return err
	}
	if err := genCode(node.rightNode); err != nil {
		return err
	}

	// 左辺と右辺にポインタがあるか確認
	existsPointer := 0
	// 0 なら存在しない、1なら左辺、2なら右辺、3は両辺 (Error)
	if node.leftNode.nodeKind == ND_LVAR {
		if node.leftNode.value.(variable).kind == TypePtr {
			existsPointer += 1
		}
	}

	if node.rightNode.nodeKind == ND_LVAR {
		if node.rightNode.value.(variable).kind == TypePtr {
			existsPointer += 2
		}
	}

	if existsPointer == 3 {
		return fmt.Errorf("genCode: not permitted of binary pointer operation")
	}

	// ADD と SUB はポインタ加減算があるので先に処理

	if node.nodeKind == ND_ADD || node.nodeKind == ND_SUB {
		generatingCode = append(generatingCode, "pop rdi\n") // right node
		generatingCode = append(generatingCode, "pop rax\n") // left node
		if existsPointer == 2 {
			generatingCode = append(generatingCode, "imul rdi, 8\n") // right node
		}
		if existsPointer == 1 {
			generatingCode = append(generatingCode, "imul rax, 8\n") // left node
		}

		if node.nodeKind == ND_ADD {
			generatingCode = append(generatingCode, "add rax, rdi\n")
		} else {
			// ND_SUB
			generatingCode = append(generatingCode, "sub rax, rdi\n")
		}
		generatingCode = append(generatingCode, "push rax\n")
		return nil
	}

	generatingCode = append(generatingCode, "pop rdi\n")
	generatingCode = append(generatingCode, "pop rax\n")

	switch node.nodeKind {
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
	return nil
}

func GenAssembleMain(nodes []*abstSyntaxNode) (codes []string, err error) {
	genInitCode()
	for _, node := range nodes {
		if err := genCode(node); err != nil {
			return nil, err
		}

		// 各式の計算結果をスタックからraxにpop
		if NoMain {
			generatingCode = append(generatingCode, "pop rax\n")
		}
	}
	genEndCode()
	return generatingCode, nil
}
