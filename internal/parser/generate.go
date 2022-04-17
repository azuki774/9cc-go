package parser

import "fmt"

var argsRegisterName = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func (cm *CodeManager) genInitCode() {
	cm.AddCode(".intel_syntax noprefix")
	cm.AddCode(".globl main")
	if cm.noMain {
		cm.AddCode("main:")
		cm.genCodePrologue("main", 0)
	}
}

func (cm *CodeManager) genEndCode() {
	// 関数を呼び出す前にスタックの状態を戻す
	cm.AddCode("mov rsp, rbp")
	cm.AddCode("pop rbp")

	// raw に入っている値を return する
	cm.AddCode("ret")
}

func (cm *CodeManager) genLocalVar(node *abstSyntaxNode) (err error) {
	// スタックの最後尾に変数のあるアドレスを入れる

	if node.nodeKind == ND_LVAR {
		// 変数の値が入っているところにポインタを移動
		cm.AddCode("mov rax, rbp")
		cm.AddCode("sub rax, %v", node.value.(variable).offset) // offset 分だけずらす

		cm.AddCode("push rax")
		// offset 分だけずらす

		return
	}

	if node.nodeKind == ND_DEREF {
		// 左辺値がDEREF(*)のとき
		cm.genCode(node.leftNode)
		return
	}

	return fmt.Errorf("genLocalVar : left value is not variable, actual = %s", node.nodeKind)
}

func (cm *CodeManager) genCodePrologue(funcName string, argsNum int) {
	cm.AddCode("push rbp")
	cm.AddCode("mov rbp, rsp") // rbp のアドレス = rsp のアドレス
	cm.AddCode("sub rsp, 256") // ローカル変数用に容量確保 32 * 8

	if funcName == "main" {
		return
	}

	if argsNum > 0 {
		cm.AddCode("mov rsp, rbp") // 一旦変数領域のスタックのベースに移動
		for i := 0; i < argsNum; i++ {
			// 受け取った引数に変数を書き換えていく
			cm.AddCode("sub rsp, 8")                         // i個目の引数のアドレスに移動
			cm.AddCode("mov [rsp], %v", argsRegisterName[i]) // i個目の引数を書き換え
		}
		cm.AddCode("mov rsp, rbp") // rbp のアドレス = rsp のアドレス
		cm.AddCode("sub rsp, 256") // ローカル変数用に容量確保 32 * 8
	}
}

func (cm *CodeManager) genCode(node *abstSyntaxNode) (err error) {
	// num | local var | = は先に処理する
	nowJumpLabel := cm.getJumpLabel()
	switch node.nodeKind {
	case ND_NIL:
		return
	case ND_NUM:
		cm.AddCode("push %d", node.value.(int))
		return
	case ND_LVAR: // local var
		if err := cm.genLocalVar(node); err != nil { // スタックの最後尾に変数のあるアドレスが入る
			return err
		}

		cm.AddCode("pop rax")        // 変数のあるアドレスがスタックから消え、raxに入る
		cm.AddCode("mov rax, [rax]") // rax の中身に rax を書きかえ、変数の値になる
		cm.AddCode("push rax")       // 変数の値をスタックに入れる
		return
	case ND_EQ:
		// 左辺をローカル変数として評価する
		if err := cm.genLocalVar(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}
		if err := cm.genCode(node.rightNode); err != nil { // スタックに右辺を計算した結果を入れるコード
			return err
		}

		cm.AddCode("pop rdi")
		cm.AddCode("pop rax")
		cm.AddCode("mov [rax], rdi") // 変数の値を直接右辺に書き換える
		cm.AddCode("push rdi\n")

		return
	case ND_ADDR: // &x
		if err := cm.genLocalVar(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}
		return
	case ND_DEREF: // *x
		if err := cm.genCode(node.leftNode); err != nil { // スタックに左辺の変数のアドレスを入れるコード
			return err
		}

		cm.AddCode("pop rax")
		cm.AddCode("mov rax, [rax]")
		cm.AddCode("push rax")

		return
	case ND_RETURN:
		if err := cm.genCode(node.leftNode); err != nil { // return する値を評価するコード
			return err
		}
		cm.AddCode("pop rax")
		// スタックを関数呼び出し前に戻す
		cm.AddCode("mov rsp, rbp")
		cm.AddCode("pop rbp")

		cm.AddCode("ret")
		return
	case ND_IF:
		// if (A) B
		cm.AddJumpLabel(1)
		if err := cm.genCode(node.leftNode); err != nil { // A
			return err
		}
		cm.AddCode("pop rax")
		cm.AddCode("cmp rax, 0")
		cm.AddCode("je  .Lend%d", nowJumpLabel)
		if err := cm.genCode(node.rightNode); err != nil { // B
			return err
		}
		cm.AddCode(".Lend%d:", nowJumpLabel)
		return
	case ND_IFELSE:
		// if (A) B else C
		cm.AddJumpLabel(1)
		if err := cm.genCode(node.leftNode); err != nil { // A
			return err
		}
		cm.AddCode("pop rax")
		cm.AddCode("cmp rax, 0")
		cm.AddCode("je  .Lelse%d", nowJumpLabel)

		if err := cm.genCode(node.rightNode.leftNode); err != nil { // B
			return err
		}

		cm.AddCode("jmp .Lend%d", nowJumpLabel)
		cm.AddCode(".Lelse%d:", nowJumpLabel)
		if err := cm.genCode(node.rightNode.rightNode); err != nil { // C
			return err
		}
		cm.AddCode(".Lend%d:", nowJumpLabel)
		return
	case ND_WHILE:
		// while (A) B
		cm.AddJumpLabel(1)
		cm.AddCode(".Lbegin%d:", nowJumpLabel)
		if err := cm.genCode(node.leftNode); err != nil { // A
			return err
		}
		cm.AddCode("pop rax")
		cm.AddCode("cmp rax, 0")
		cm.AddCode("je  .Lend%d", nowJumpLabel)
		if err := cm.genCode(node.rightNode); err != nil { // B
			return err
		}
		cm.AddCode("jmp .Lbegin%d", nowJumpLabel)
		cm.AddCode(".Lend%d:", nowJumpLabel)
		return
	case ND_FOR:
		// for (A; B; C) D
		cm.AddJumpLabel(1)
		A := node.leftNode.leftNode
		B := node.leftNode.rightNode
		C := node.rightNode.leftNode
		D := node.rightNode.rightNode

		if err := cm.genCode(A); err != nil {
			return err
		}
		cm.AddCode(".Lbegin%d:", nowJumpLabel)
		if err := cm.genCode(B); err != nil {
			return err
		}
		cm.AddCode("pop rax")
		cm.AddCode("cmp rax, 0")
		cm.AddCode("je  .Lend%d", nowJumpLabel)
		if err := cm.genCode(D); err != nil {
			return err
		}
		if err := cm.genCode(C); err != nil {
			return err
		}
		cm.AddCode("jmp .Lbegin%d", nowJumpLabel)
		cm.AddCode(".Lend%d:", nowJumpLabel)
		return
	case ND_BLOCK:
		// { stmt* }
		stmtNodeList := node.value.([]*abstSyntaxNode)
		for _, nowNode := range stmtNodeList {
			if err := cm.genCode(nowNode); err != nil {
				return err
			}
		}
		return
	case ND_FUNDEF:
		funcName := node.value.(string)
		cm.AddCode("%s:", funcName)

		// この変数が何変数関数か調べる
		var argsNum int = 0
		argsNode := node.leftNode
		if argsNode.value != nil {
			argsNum = len(argsNode.value.([]*abstSyntaxNode))
		}

		cm.genCodePrologue(funcName, argsNum)
		if err := cm.genCode(node.rightNode); err != nil {
			return err
		}
		return
	case ND_FUNDEF_ARGS:
		argsNodes := node.value.([]*abstSyntaxNode)
		for _, v := range argsNodes {
			if err := cm.genCode(v); err != nil { // 変数を定義する
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
			if err := cm.genCode(v); err != nil {
				return err
			}
			cm.AddCode("pop %s", argsRegisterName[i])
		}
		// rsp を 16の倍数にする調整
		cm.AddJumpLabel(1)
		cm.AddCode("mov rax, rsp")
		cm.AddCode("and rax, 15")                // 下4bitのみにマスキング
		cm.AddCode("jnz .Lcall%d", cm.jumpLabel) // 下4bit != 0
		cm.AddCode("mov rax, 0")                 // TODO: rax は引数の数
		cm.AddCode("call %s", funcName)
		cm.AddCode("jmp .Lend%d", cm.jumpLabel)
		cm.AddCode(".Lcall%d:", cm.jumpLabel) // 16の倍数になっていなくて、8ずらすときはここから
		cm.AddCode("sub rsp, 8")
		cm.AddCode("mov rax, 0") // TODO: rax は引数の数
		cm.AddCode("call %s", funcName)
		cm.AddCode("add rsp, 8")
		cm.AddCode(".Lend%d:", cm.jumpLabel)
		cm.AddCode("push rax")
		return
	}

	// 二項演算系

	if err := cm.genCode(node.leftNode); err != nil {
		return err
	}
	if err := cm.genCode(node.rightNode); err != nil {
		return err
	}

	// 左辺と右辺にポインタがあるか確認
	existsPointer := 0
	// 0 なら存在しない、1なら左辺、2なら右辺、3は両辺 (Error)
	if node.leftNode.nodeKind == ND_LVAR {
		if node.leftNode.value.(variable).kind.primKind == TypePtr {
			existsPointer += 1
		}
	}

	if node.rightNode.nodeKind == ND_LVAR {
		if node.rightNode.value.(variable).kind.primKind == TypePtr {
			existsPointer += 2
		}
	}

	if existsPointer == 3 {
		return fmt.Errorf("genCode: not permitted of binary pointer operation")
	}

	// ADD と SUB はポインタ加減算があるので先に処理

	if node.nodeKind == ND_ADD || node.nodeKind == ND_SUB {
		cm.AddCode("pop rdi") // right node
		cm.AddCode("pop rax") // left node
		if existsPointer == 2 {
			cm.AddCode("imul rdi, 8") // right node
		}
		if existsPointer == 1 {
			cm.AddCode("imul rax, 8") // left node
		}

		if node.nodeKind == ND_ADD {
			cm.AddCode("add rax, rdi") // right node
		} else {
			// ND_SUB
			cm.AddCode("sub rax, rdi") // right node
		}
		cm.AddCode("push rax") // right node
		return nil
	}

	cm.AddCode("pop rdi")
	cm.AddCode("pop rax")

	switch node.nodeKind {
	case ND_MUL:
		cm.AddCode("imul rax, rdi")
	case ND_DIV:
		cm.AddCode("cqo")
		cm.AddCode("idiv rdi")
	case ND_COMP:
		cm.AddCode("cmp rax, rdi")
		cm.AddCode("sete al")
		cm.AddCode("movzb rax, al")
	case ND_NOTEQ:
		cm.AddCode("cmp rax, rdi")
		cm.AddCode("setne al")
		cm.AddCode("movzb rax, al")
	case ND_LT:
		cm.AddCode("cmp rax, rdi")
		cm.AddCode("setl al")
		cm.AddCode("movzb rax, al")
	case ND_LTQ:
		cm.AddCode("cmp rax, rdi")
		cm.AddCode("setle al")
		cm.AddCode("movzb rax, al")
	}

	cm.AddCode("push rax")
	return nil
}
