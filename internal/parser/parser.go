package parser

var (
	ts *tokenStream
)

var (
	ND_UNDEFINED = 0
	ND_ADD       = 1
	ND_SUB       = 2
	ND_MUL       = 3
	ND_DIV       = 4
	ND_NUM       = 10
)

type abstSyntaxNode struct {
	nodeKind  int
	leftNode  *abstSyntaxNode
	rightNode *abstSyntaxNode
	value     interface{}
}

func makeNewAbstSyntaxNode(nodeKind int, leftNode *abstSyntaxNode, rightNode *abstSyntaxNode, value interface{}) *abstSyntaxNode {
	return &abstSyntaxNode{nodeKind: nodeKind, leftNode: leftNode, rightNode: rightNode, value: value}
}

func Expr_expr() (node *abstSyntaxNode) {
	// expr    = mul ("+" mul | "-" mul )*
	node = Expr_mul()
	// + か - があるとき
	for {
		nToken := ts.nextPeekToken()
		if nToken.kind == TK_SYMBOL {
			b := nToken.value.(byte)
			switch b {
			case BYTE_SYMBOL_ADD:
				ts.nextToken()
				node = makeNewAbstSyntaxNode(ND_ADD, node, Expr_mul(), nil)
				continue
			case BYTE_SYMBOL_SUB:
				ts.nextToken()
				node = makeNewAbstSyntaxNode(ND_SUB, node, Expr_mul(), nil)
				continue
			}
		}

		// + Token でも - Token でもないとき
		break
	}

	return node
}

func Expr_mul() (node *abstSyntaxNode) {
	// mul     = primary ("*" primary | "/" primary )*
	node = Expr_primary()
	// * か / があるとき
	for {
		nToken := ts.nextPeekToken()
		if nToken.kind == TK_SYMBOL {
			b := nToken.value.(byte)
			switch b {
			case BYTE_SYMBOL_MUL:
				ts.nextToken()
				node = makeNewAbstSyntaxNode(ND_MUL, node, Expr_primary(), nil)
				continue
			case BYTE_SYMBOL_DIV:
				ts.nextToken()
				node = makeNewAbstSyntaxNode(ND_DIV, node, Expr_primary(), nil)
				continue
			}
		}

		// * Token でも / Token でもないとき
		break
	}

	return node
}

func Expr_primary() (node *abstSyntaxNode) {
	// primary = num | "(" expr ")"
	nToken := ts.nextPeekToken()
	if nToken.kind == TK_SYMBOL {
		b := nToken.value.(byte)
		if b == BYTE_LEFT_PAT {
			ts.nextToken() // (
			node = Expr_expr()
			ts.nextToken() // )
			return node
		}

	}
	node = makeNewAbstSyntaxNode(ND_NUM, nil, nil, ts.nextToken().value) // Token は 1つ進む
	// Tokenは1つ進む
	return node
}

func ParserMain(tokens []Token) (topNode *abstSyntaxNode) {
	ts = newTokenStream(tokens)
	topNode = Expr_expr()
	return topNode
}
