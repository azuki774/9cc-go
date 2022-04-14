package parser

import "fmt"

var (
	ts *tokenStream
)

var (
	ND_UNDEFINED    = 0
	ND_NIL          = 1 // nil (no code)
	ND_NUM          = 10
	ND_ADD          = 11
	ND_SUB          = 12
	ND_MUL          = 13
	ND_DIV          = 14
	ND_COMP         = 21 // ==
	ND_NOTEQ        = 22 // !=
	ND_LT           = 23 // <
	ND_LTQ          = 24 // <=
	ND_EQ           = 25 // =
	ND_ADDR         = 26 // &hoge
	ND_DEREF        = 27 // *hoge
	ND_LVAR         = 31 // local variable, value に struct Var
	ND_RETURN       = 41
	ND_IF           = 42
	ND_ELSE         = 43
	ND_IFELSE       = 44 // elseありのIF
	ND_WHILE        = 45
	ND_FOR          = 46
	ND_BLOCK        = 47 // { stmt* } : value に stmt* に含まれるabstSyntaxNode のスライス
	ND_FUNCALL      = 48 // value に呼び出す関数名、leftNode に ND_FUNCALL_ARG
	ND_FUNCALL_ARGS = 49 // value に引数たちの abstSyntaxNode のスライス
	ND_FUNDEF       = 50 // value に関数名、leftNode に ND_FUNDEF_ARGS, rightNode に 関数のstmt
	ND_FUNDEF_ARGS  = 51 // value に args の node のスライスを詰める
)

type TypeKind string

const (
	TypeInt = TypeKind("int")
	TypePtr = TypeKind("pointer")
)

type variableManager struct {
	varList    map[string]variable
	nextoffset int
}

type variable struct {
	kind   TypeKind
	ptrTo  *variable
	offset int
}

type abstSyntaxNode struct {
	nodeKind  int
	leftNode  *abstSyntaxNode
	rightNode *abstSyntaxNode
	value     interface{} // num の値や、local variable の offset を入れる
}

func makeNewAbstSyntaxNode(nodeKind int, leftNode *abstSyntaxNode, rightNode *abstSyntaxNode, value interface{}) *abstSyntaxNode {
	return &abstSyntaxNode{nodeKind: nodeKind, leftNode: leftNode, rightNode: rightNode, value: value}
}

func makeNewVariableManager() *variableManager {
	return &variableManager{varList: map[string]variable{}, nextoffset: 8}
}

func (v *variableManager) reset() {
	v.varList = map[string]variable{}
	v.nextoffset = 8
}

// 変数をvariableManager
func (v *variableManager) add(varname string, kind TypeKind) (nvar variable, err error) {
	if _, ok := v.varList[varname]; ok {
		// 変数が定義済
		return variable{}, fmt.Errorf("%s is already defined", varname)
	} else {
		// 変数が未定義 -> 追加
		switch kind {
		case TypeInt:
			nvar = variable{kind: TypeInt, ptrTo: nil, offset: v.nextoffset}
			v.nextoffset += 8
		case TypePtr:
			nvar = variable{kind: TypePtr, ptrTo: nil, offset: v.nextoffset}
			v.nextoffset += 8
		}

	}

	v.varList[varname] = nvar
	return nvar, nil
}

func ParserMain(tokens []Token) (nodes []*abstSyntaxNode, err error) {
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
