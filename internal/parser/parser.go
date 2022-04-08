package parser

var (
	ts *tokenStream
)

var (
	ND_UNDEFINED   = 0
	ND_NIL         = 1 // nil (no code)
	ND_NUM         = 10
	ND_ADD         = 11
	ND_SUB         = 12
	ND_MUL         = 13
	ND_DIV         = 14
	ND_COMP        = 21 // ==
	ND_NOTEQ       = 22 // !=
	ND_LT          = 23 // <
	ND_LTQ         = 24 // <=
	ND_EQ          = 25 // =
	ND_LVAR        = 31 // local variable
	ND_RETURN      = 41
	ND_IF          = 42
	ND_ELSE        = 43
	ND_IFELSE      = 44 // elseありのIF
	ND_WHILE       = 45
	ND_FOR         = 46
	ND_BLOCK       = 47 // { stmt* } : value に stmt* に含まれるabstSyntaxNode のスライス
	ND_FUNCALL     = 48 // value に呼び出す関数名、leftNode に引数のnode
	ND_FUNCALL_ARG = 49 // value に引数たちの abstSyntaxNode のスライス
	ND_FUNDEF      = 50 // value に関数名、leftNode に引数のnode, rightNode に 関数のstmt
)

type abstSyntaxNode struct {
	nodeKind  int
	leftNode  *abstSyntaxNode
	rightNode *abstSyntaxNode
	value     interface{} // num の値や、local variable の offset を入れる
}

func makeNewAbstSyntaxNode(nodeKind int, leftNode *abstSyntaxNode, rightNode *abstSyntaxNode, value interface{}) *abstSyntaxNode {
	return &abstSyntaxNode{nodeKind: nodeKind, leftNode: leftNode, rightNode: rightNode, value: value}
}

func ParserMain(tokens []Token) (nodes []*abstSyntaxNode, err error) {
	localVar = map[string]int{}
	ts = newTokenStream(tokens)
	if !NoMain {
		nodes, err = Expr_program(ts)
	} else {
		nodes, err = Expr_programNoMain(ts)
	}

	if err != nil {
		return nil, err
	}

	return nodes, nil
}
