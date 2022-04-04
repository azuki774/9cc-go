package parser

var (
	ts *tokenStream
)

var (
	ND_UNDEFINED = 0
	ND_NUM       = 1
	ND_ADD       = 11
	ND_SUB       = 12
	ND_MUL       = 13
	ND_DIV       = 14
	ND_COMP      = 21 // ==
	ND_NOTEQ     = 22 // !=
	ND_LT        = 23 // <
	ND_LTQ       = 24 // <=
	ND_EQ        = 25 // =
	ND_LVAR      = 31 // local variable
	ND_RETURN    = 41
	ND_IF        = 42
	ND_ELSE      = 43
	ND_IFELSE    = 44 // elseありのIF
)

type abstSyntaxNode struct {
	nodeKind  int
	leftNode  *abstSyntaxNode
	rightNode *abstSyntaxNode
	value     interface{} // num の値や、local variable の offset を入れる
}

var localVar map[string]int // varName -> offset

func makeNewAbstSyntaxNode(nodeKind int, leftNode *abstSyntaxNode, rightNode *abstSyntaxNode, value interface{}) *abstSyntaxNode {
	return &abstSyntaxNode{nodeKind: nodeKind, leftNode: leftNode, rightNode: rightNode, value: value}
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
		ts.nextToken()                                               // if
		err = ts.nextExpectReadToken(Token{kind: TK_SYMBOL_LEFTPAT}) // (
		if err != nil {
			return nil, err
		}

		eA, err := Expr_expr()
		if err != nil {
			return nil, err
		}

		err = ts.nextExpectReadToken(Token{kind: TK_SYMBOL_RIGHTPAT}) // )
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
		switch nToken.kind {
		case TK_SYMBOL_ADD:
			ts.nextToken()
			e, err := Expr_mul()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_ADD, node, e, nil)
			continue
		case TK_SYMBOL_SUB:
			ts.nextToken()
			e, err := Expr_mul()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_SUB, node, e, nil)
			continue
		}

		// + Token でも - Token でもないとき
		break
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
		case TK_SYMBOL_MUL:
			ts.nextToken()
			e, err := Expr_unary()
			if err != nil {
				return nil, err
			}
			node = makeNewAbstSyntaxNode(ND_MUL, node, e, nil)
			continue
		case TK_SYMBOL_DIV:
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
	// +x -> x, -x -> 0-x
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_SYMBOL_ADD:
		ts.nextToken() // + を読み飛ばす
		node, err = Expr_primary()
		if err != nil {
			return nil, err
		}
		return node, nil
	case TK_SYMBOL_SUB:
		ts.nextToken()
		leftNode := makeNewAbstSyntaxNode(ND_NUM, nil, nil, 0)
		e, err := Expr_unary()
		if err != nil {
			return nil, err
		}
		node = makeNewAbstSyntaxNode(ND_SUB, leftNode, e, nil) // "0 -" x に対応
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
	case TK_SYMBOL_LEFTPAT: // ( expr )
		ts.nextToken() // (
		node, err = Expr_expr()
		if err != nil {
			return nil, err
		}
		err = ts.nextExpectReadToken(Token{kind: TK_SYMBOL_RIGHTPAT}) // )
		if err != nil {
			return nil, err
		}
	case TK_NUM:
		node = makeNewAbstSyntaxNode(ND_NUM, nil, nil, ts.nextToken().value) // Token は 1つ進む
	case TK_IDENT:
		name := nToken.value.(string)
		if offset, ok := localVar[name]; ok {
			// 変数が定義済
			node = makeNewAbstSyntaxNode(ND_LVAR, nil, nil, int(offset))
		} else {
			// 変数が初出
			offset = (len(localVar) + 1) * 8
			localVar[name] = offset
			node = makeNewAbstSyntaxNode(ND_LVAR, nil, nil, int(offset))
		}
		ts.nextToken() // 変数トークンを消化する
	}

	return node, nil
}

func ParserMain(tokens []Token) (topNodes []*abstSyntaxNode, err error) {
	localVar = map[string]int{}
	ts = newTokenStream(tokens)
	for {
		if !ts.ok() {
			break
		}

		topNode, err := Expr_stmt()
		if err != nil {
			return nil, err
		}
		topNodes = append(topNodes, topNode)
	}
	return topNodes, nil
}
