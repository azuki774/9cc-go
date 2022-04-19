package parser

import "fmt"

// var localVar map[string]variable // varName -> offset
var localVar *variableManager

func Expr_program(ts *tokenStream) (nodes []*abstSyntaxNode, err error) {
	localVar = makeNewVariableManager()

	for {
		if !ts.ok() {
			break
		}

		localVar.reset()

		// ident "(" ")" stmt
		if ts.nextPeekToken().kind == TK_IDENT {
			ts.nextToken() // int

			if ts.nextPeekToken().kind != TK_IDENT {
				return nil, fmt.Errorf("Expr_program : not found the function name")
			}

			funcName := ts.nextPeekToken().value.(string)
			ts.nextToken() // ident
			ts.nextToken() // (

			argsNode := makeNewAbstSyntaxNode(ND_FUNDEF_ARGS, nil, nil, []*abstSyntaxNode{})

			for {
				if ts.nextPeekToken().kind == TK_COMMA {
					ts.nextToken() // ,
				}
				if ts.nextPeekToken().kind == TK_RIGHTPAT {
					break
				}
				newVarNode, err := Expr_primary() // 変数定義node
				if err != nil {
					return nil, err
				}
				argsNode.value = append(argsNode.value.([]*abstSyntaxNode), newVarNode)
			}

			ts.nextToken() // )
			stmtNode, err := Expr_stmt()
			if err != nil {
				return nil, err
			}
			topNode := makeNewAbstSyntaxNode(ND_FUNDEF, argsNode, stmtNode, funcName)
			nodes = append(nodes, topNode)
		} else {
			return nil, fmt.Errorf("Expr_program : not found definition of function")
		}

	}

	return nodes, nil
}

func Expr_programNoMain(ts *tokenStream) (nodes []*abstSyntaxNode, err error) {
	localVar = makeNewVariableManager()

	for {
		if !ts.ok() {
			break
		}

		topNode, err := Expr_stmt()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, topNode)
	}
	return nodes, nil
}

func Expr_stmt() (node *abstSyntaxNode, err error) {
	// stmt = expr ";" | "return" expr ";" | "if" "(" expr ")" stmt ("else" stmt)?
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_RETURN:
		// return a
		ts.nextToken() // return
		e, err := Expr_expr()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_RETURN, e, nil, nil)
	case TK_IF:
		ts.nextToken()                                        // if
		err = ts.nextExpectReadToken(Token{kind: TK_LEFTPAT}) // (
		if err != nil {
			return nil, err
		}

		eA, err := Expr_expr()
		if err != nil {
			return nil, err
		}

		err = ts.nextExpectReadToken(Token{kind: TK_RIGHTPAT}) // )
		if err != nil {
			return nil, err
		}

		eB, err := Expr_stmt()
		if err != nil {
			return nil, err
		}

		if ts.nextPeekToken().kind == TK_ELSE {
			ts.nextToken() // else
			eC, err := Expr_stmt()
			if err != nil {
				return nil, err
			}
			nodeSuc := makeNewAbstSyntaxNode(ND_ELSE, eB, eC, nil)
			node = makeNewAbstSyntaxNode(ND_IFELSE, eA, nodeSuc, nil)
		} else {
			node = makeNewAbstSyntaxNode(ND_IF, eA, eB, nil)
		}

		return node, nil
	case TK_WHILE:
		// "while" "(" expr ")" stmt
		ts.nextToken() // while
		err = ts.nextExpectReadToken(Token{kind: TK_LEFTPAT})
		if err != nil {
			return nil, err
		}

		eA, err := Expr_expr()
		if err != nil {
			return nil, err
		}

		err = ts.nextExpectReadToken(Token{kind: TK_RIGHTPAT})
		if err != nil {
			return nil, err
		}

		eB, err := Expr_stmt()
		if err != nil {
			return nil, err
		}

		node = makeNewAbstSyntaxNode(ND_WHILE, eA, eB, nil)
		return node, nil
	case TK_FOR:
		// "for" "(" expr? ";" expr? ";" expr? ")" stmt
		// for (A; B; C) D

		ts.nextToken() // for
		err = ts.nextExpectReadToken(Token{kind: TK_LEFTPAT})
		if err != nil {
			return nil, err
		}

		var eA *abstSyntaxNode = makeNewAbstSyntaxNode(ND_NIL, nil, nil, nil)
		var eB *abstSyntaxNode = makeNewAbstSyntaxNode(ND_NIL, nil, nil, nil)
		var eC *abstSyntaxNode = makeNewAbstSyntaxNode(ND_NIL, nil, nil, nil)

		// A
		if ts.nextPeekToken().kind != TK_SEMICOLON {
			eA, err = Expr_expr()
			if err != nil {
				return nil, err
			}
		}
		err = ts.nextExpectReadToken(Token{kind: TK_SEMICOLON})
		if err != nil {
			return nil, err
		}

		// B
		if ts.nextPeekToken().kind != TK_SEMICOLON {
			eB, err = Expr_expr()
			if err != nil {
				return nil, err
			}
		}
		err = ts.nextExpectReadToken(Token{kind: TK_SEMICOLON})
		if err != nil {
			return nil, err
		}

		// C
		if ts.nextPeekToken().kind != TK_RIGHTPAT {
			eC, err = Expr_expr()
			if err != nil {
				return nil, err
			}
		}

		err = ts.nextExpectReadToken(Token{kind: TK_RIGHTPAT})
		if err != nil {
			return nil, err
		}
		eD, err := Expr_stmt()
		if err != nil {
			return nil, err
		}

		// for (A; B; C) D
		eAB := makeNewAbstSyntaxNode(ND_UNDEFINED, eA, eB, nil)
		eCD := makeNewAbstSyntaxNode(ND_UNDEFINED, eC, eD, nil)
		node = makeNewAbstSyntaxNode(ND_FOR, eAB, eCD, nil)

		return node, nil

	case TK_BLOCKL:
		ts.nextToken() // {
		var blockNDList []*abstSyntaxNode
		for { // stmt*
			if ts.nextPeekToken().kind == TK_BLOCKR {
				ts.nextToken() // }
				break
			}
			newNode, err := Expr_stmt()
			if err != nil {
				return nil, err
			}
			blockNDList = append(blockNDList, newNode)
		}

		node = makeNewAbstSyntaxNode(ND_BLOCK, nil, nil, blockNDList)
		return node, nil

	default:
		node, err = Expr_expr()
		if err != nil {
			return nil, err
		}
	}

	err = ts.nextExpectReadToken(Token{kind: TK_SEMICOLON}) // ;
	if err != nil {
		return nil, err
	}
	return node, nil
}

func Expr_expr() (node *abstSyntaxNode, err error) {
	// expr = assign
	node, err = Expr_assign()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func Expr_assign() (node *abstSyntaxNode, err error) {
	// assign = equality ("=" assign)?
	node, err = Expr_equality()
	if err != nil {
		return nil, err
	}

	nToken := ts.nextPeekToken()
	if nToken.kind == TK_EQ {
		ts.nextToken() // =
		e, err := Expr_assign()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_EQ, node, e, nil)
	}

	return node, nil
}

func Expr_equality() (node *abstSyntaxNode, err error) {
	// equality = relational ("==" relational | "!=" relational)*
	node, err = Expr_relational()
	if err != nil {
		return nil, err
	}

	// == か != があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_COMP:
			ts.nextToken()
			e, err := Expr_relational()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_COMP, node, e, nil)
			continue
		case TK_NOTEQ:
			ts.nextToken()
			e, err := Expr_relational()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_NOTEQ, node, e, nil)
			continue
		}

		// == Token でも != でもないとき
		break
	}

	return node, nil
}

func Expr_relational() (node *abstSyntaxNode, err error) {
	// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
	node, err = Expr_add()
	if err != nil {
		return nil, err
	}

	// <, > か <=, >= があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_LT:
			ts.nextToken()
			e, err := Expr_add()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_LT, node, e, nil)
			continue
		case TK_LTQ:
			ts.nextToken()
			e, err := Expr_add()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_LTQ, node, e, nil)
			continue
		case TK_GT:
			// TK_LT の左右入れ替え
			ts.nextToken()
			e, err := Expr_add()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_LT, e, node, nil)
			continue
		case TK_GTQ:
			// TK_LTQ の左右入れ替え
			ts.nextToken()
			e, err := Expr_add()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_LTQ, e, node, nil)
			continue
		}

		// == Token でも != でもないとき
		break
	}

	return node, nil
}

func Expr_add() (node *abstSyntaxNode, err error) {
	// add = mul ("+" mul | "-" mul )*
	node, err = Expr_mul()
	if err != nil {
		return nil, err
	}

	// + か - があるとき
	for {
		nToken := ts.nextPeekToken()
		if nToken.kind == TK_ADD || nToken.kind == TK_SUB {
			// node + e or node - e
			ts.nextToken() // + or -
			e, err := Expr_mul()
			if err != nil {
				return nil, err
			}

			if nToken.kind == TK_ADD {
				node = makeNewAbstSyntaxNode(ND_ADD, node, e, nil)
			} else if nToken.kind == TK_SUB {
				node = makeNewAbstSyntaxNode(ND_SUB, node, e, nil)
			}
		} else {
			// + Token でも - Token でもないとき
			break
		}
	}

	return node, nil
}

func Expr_mul() (node *abstSyntaxNode, err error) {
	// mul     = unary ("*" unary | "/" unary)*
	node, err = Expr_unary()
	if err != nil {
		return nil, err
	}
	// * か / があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_MUL:
			ts.nextToken()
			e, err := Expr_unary()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_MUL, node, e, nil)
			continue
		case TK_DIV:
			ts.nextToken()
			e, err := Expr_unary()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_DIV, node, e, nil)
			continue
		}

		// * Token でも / Token でもないとき
		break
	}

	return node, nil
}

func Expr_unary() (node *abstSyntaxNode, err error) {
	// unary   = ("+" | "-")? primary
	// "*" unary
	// "&" unary
	// +x -> x, -x -> 0-x
	// "sizeof" unary
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_ADD:
		ts.nextToken() // + を読み飛ばす
		node, err = Expr_primary()
		if err != nil {
			return nil, err
		}
		return node, nil
	case TK_SUB:
		ts.nextToken()
		leftNode := makeNewAbstSyntaxNode(ND_NUM, nil, nil, 0)
		e, err := Expr_primary()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_SUB, leftNode, e, nil) // "0 -" x に対応
		return node, nil
	case TK_MUL: // DEREF
		ts.nextToken()
		e, err := Expr_unary()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_DEREF, e, nil, nil) // *e
		return node, nil
	case TK_AND: // ADDR
		ts.nextToken()
		e, err := Expr_unary()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_ADDR, e, nil, nil) // &e
		return node, nil
	case TK_SIZEOF:
		ts.nextToken()
		e, err := Expr_unary()
		if err != nil {
			return nil, err
		}

		// e を見て sizeof e を評価してしまう
		size, err := getSizeOf(e)
		node = makeNewAbstSyntaxNode(ND_NUM, nil, nil, size)
		return node, nil
	}

	// + も - もないとき
	node, err = Expr_primary()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func Expr_primary() (node *abstSyntaxNode, err error) {
	// primary = num | ident | "(" expr ")"
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_LEFTPAT: // ( expr )
		ts.nextToken() // (
		node, err = Expr_expr()
		if err != nil {
			return nil, err
		}
		err = ts.nextExpectReadToken(Token{kind: TK_RIGHTPAT}) // )
		if err != nil {
			return nil, err
		}
	case TK_NUM:
		node = makeNewAbstSyntaxNode(ND_NUM, nil, nil, ts.nextToken().value) // Token は 1つ進む
	case TK_IDENT:
		name := nToken.value.(string)

		if name == "int" {
			// 変数定義
			node, err = expr_defVar()
			if err != nil {
				return nil, err
			}
			return node, nil
		}

		if ts.nextPeekToken().kind == TK_MUL {
			// 変数評価 (*type)
			node, err = expr_evalVar()
			if err != nil {
				return nil, err
			}
			return node, nil
		}

		ts.nextToken() // 一旦変数名を読み捨てる

		if ts.nextPeekToken().kind != TK_LEFTPAT {
			// 変数評価
			ts.backToken()
			node, err = expr_evalVar()
			if err != nil {
				return nil, err
			}
			return node, nil
		} else {
			// 関数評価
			ts.backToken()
			node, err = expr_evalFunc()
			if err != nil {
				return nil, err
			}
			return node, nil
		}

	}

	return node, nil
}

func expr_defVar() (node *abstSyntaxNode, err error) {
	ts.nextToken() // int
	if ts.nextPeekToken().kind != TK_IDENT && ts.nextPeekToken().kind != TK_MUL {
		return nil, fmt.Errorf("Expr_primary: variable name error")
	}

	// ポインタの*のTokenの数を計算
	pointerCount := 0
	for {
		if ts.nextPeekToken().kind == TK_MUL {
			ts.nextToken() // *
			pointerCount += 1
		} else {
			break
		}
	}

	varName := ts.nextPeekToken().value.(string)
	ts.nextToken() // 変数名

	var nvar variable
	nvar, err = localVar.add(varName, makeTypeKind(TypeInt, pointerCount))
	// if pointerCount == 0 {
	// 	nvar, err = localVar.add(varName, TypeInt)
	// } else {
	// 	nvar, err = localVar.add(varName, TypePtr)
	// }

	if err != nil {
		return nil, err
	}

	node = makeNewAbstSyntaxNode(ND_NIL, nil, nil, nvar)
	return node, nil
}

func expr_evalFunc() (node *abstSyntaxNode, err error) {
	funcName := ts.nextPeekToken().value.(string)
	ts.nextToken() // 関数名
	ts.nextToken() // (

	varNode := makeNewAbstSyntaxNode(ND_FUNCALL_ARGS, nil, nil, []*abstSyntaxNode{})
	for {
		if ts.nextPeekToken().kind == TK_COMMA {
			ts.nextToken() // ,
		}
		if ts.nextPeekToken().kind == TK_RIGHTPAT {
			break
		}
		newVarNode, err := Expr_add() // 変数定義node
		if err != nil {
			return nil, err
		}
		varNode.value = append(varNode.value.([]*abstSyntaxNode), newVarNode)
	}

	ts.nextToken() // )
	node = makeNewAbstSyntaxNode(ND_FUNCALL, varNode, nil, funcName)
	return node, nil
}

func expr_evalVar() (node *abstSyntaxNode, err error) {
	varName := ts.nextPeekToken().value.(string)
	ts.nextToken() // 変数名

	if nvar, ok := localVar.varList[varName]; ok {
		// 変数が定義済
		node = makeNewAbstSyntaxNode(ND_LVAR, nil, nil, nvar)
	} else {
		// 変数が未定義
		return nil, fmt.Errorf("%s is not defined yet", varName)
	}
	return node, nil
}
