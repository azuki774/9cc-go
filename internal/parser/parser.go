package parser

import "fmt"

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

func Expr_stmt() (node *abstSyntaxNode) {
	// stmt = expr ";" | "return" expr ";"
	nToken := ts.nextPeekToken()
	if nToken.kind == TK_RETURN {
		// return a
		ts.nextToken() // return
		node = makeNewAbstSyntaxNode(ND_RETURN, Expr_expr(), nil, nil)
	} else {
		// expr
		node = Expr_expr()
	}

	nToken = ts.nextPeekToken() // ;
	if nToken.kind != TK_SEMICOLON {
		panic(fmt.Errorf("Expr_stmt : not found semicolon"))
	}
	ts.nextToken() // ;
	return node
}

func Expr_expr() (node *abstSyntaxNode) {
	// expr = assign
	node = Expr_assign()

	return node
}

func Expr_assign() (node *abstSyntaxNode) {
	// assign = equality ("=" assign)?
	node = Expr_equality()

	nToken := ts.nextPeekToken()
	if nToken.kind == TK_EQ {
		ts.nextToken() // =
		node = makeNewAbstSyntaxNode(ND_EQ, node, Expr_assign(), nil)
	}

	return node
}

func Expr_equality() (node *abstSyntaxNode) {
	// equality = relational ("==" relational | "!=" relational)*
	node = Expr_relational()
	// == か != があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_COMP:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_COMP, node, Expr_relational(), nil)
			continue
		case TK_NOTEQ:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_NOTEQ, node, Expr_relational(), nil)
			continue
		}

		// == Token でも != でもないとき
		break
	}

	return node
}

func Expr_relational() (node *abstSyntaxNode) {
	// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
	node = Expr_add()
	// <, > か <=, >= があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_LT:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_LT, node, Expr_add(), nil)
			continue
		case TK_LTQ:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_LTQ, node, Expr_add(), nil)
			continue
		case TK_GT:
			// TK_LT の左右入れ替え
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_LT, Expr_add(), node, nil)
			continue
		case TK_GTQ:
			// TK_LTQ の左右入れ替え
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_LTQ, Expr_add(), node, nil)
			continue
		}

		// == Token でも != でもないとき
		break
	}

	return node
}

func Expr_add() (node *abstSyntaxNode) {
	// add = mul ("+" mul | "-" mul )*
	node = Expr_mul()
	// + か - があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_SYMBOL_ADD:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_ADD, node, Expr_mul(), nil)
			continue
		case TK_SYMBOL_SUB:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_SUB, node, Expr_mul(), nil)
			continue
		}

		// + Token でも - Token でもないとき
		break
	}
	return node
}

func Expr_mul() (node *abstSyntaxNode) {
	// mul     = unary ("*" unary | "/" unary)*
	node = Expr_unary()
	// * か / があるとき
	for {
		nToken := ts.nextPeekToken()
		switch nToken.kind {
		case TK_SYMBOL_MUL:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_MUL, node, Expr_mul(), nil)
			continue
		case TK_SYMBOL_DIV:
			ts.nextToken()
			node = makeNewAbstSyntaxNode(ND_DIV, node, Expr_mul(), nil)
			continue
		}

		// * Token でも / Token でもないとき
		break
	}

	return node
}

func Expr_unary() (node *abstSyntaxNode) {
	// unary   = ("+" | "-")? primary
	// +x -> x, -x -> 0-x
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_SYMBOL_ADD:
		ts.nextToken() // + を読み飛ばす
		node = Expr_primary()
		return node
	case TK_SYMBOL_SUB:
		ts.nextToken()
		leftNode := makeNewAbstSyntaxNode(ND_NUM, nil, nil, 0)
		node = makeNewAbstSyntaxNode(ND_SUB, leftNode, Expr_primary(), nil) // "0 -" x に対応
		return node
	}

	// + も - もないとき
	node = Expr_primary()
	return node
}

func Expr_primary() (node *abstSyntaxNode) {
	// primary = num | ident | "(" expr ")"
	nToken := ts.nextPeekToken()
	switch nToken.kind {
	case TK_SYMBOL_LEFTPAT: // ( expr )
		ts.nextToken() // (
		node = Expr_expr()
		ts.nextToken() // )
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

	return node
}

func ParserMain(tokens []Token) (topNodes []*abstSyntaxNode) {
	localVar = map[string]int{}
	ts = newTokenStream(tokens)
	for {
		if !ts.ok() {
			break
		}

		topNode := Expr_stmt()
		topNodes = append(topNodes, topNode)
	}
	return topNodes
}
