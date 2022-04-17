package parser

import "fmt"

var (
	ts *tokenStream
)

type NodeKind string

const (
	ND_UNDEFINED    = NodeKind("ND_UNDEFINED")
	ND_NIL          = NodeKind("ND_NIL") // nil (no code)
	ND_NUM          = NodeKind("ND_NUM")
	ND_ADD          = NodeKind("ND_ADD")
	ND_SUB          = NodeKind("ND_SUB")
	ND_MUL          = NodeKind("ND_MUL")
	ND_DIV          = NodeKind("ND_DIV")
	ND_COMP         = NodeKind("ND_COMP")  // ==
	ND_NOTEQ        = NodeKind("ND_NOTEQ") // !=
	ND_LT           = NodeKind("ND_LT")    // <
	ND_LTQ          = NodeKind("ND_LTQ")   // <=
	ND_EQ           = NodeKind("ND_EQ")    // =
	ND_ADDR         = NodeKind("ND_ADDR")  // &hoge
	ND_DEREF        = NodeKind("ND_DEREF") // *hoge
	ND_LVAR         = NodeKind("ND_LVAR")  // local variable, value に struct Var
	ND_RETURN       = NodeKind("ND_RETURN")
	ND_IF           = NodeKind("ND_IF")
	ND_ELSE         = NodeKind("ND_ELSE")
	ND_IFELSE       = NodeKind("ND_IFELSE") // elseありのIF
	ND_WHILE        = NodeKind("ND_WHILE")
	ND_FOR          = NodeKind("ND_FOR")
	ND_BLOCK        = NodeKind("ND_BLOCK")        // { stmt* } : value に stmt* に含まれるabstSyntaxNode のスライス
	ND_FUNCALL      = NodeKind("ND_FUNCALL")      // value に呼び出す関数名、leftNode に ND_FUNCALL_ARG
	ND_FUNCALL_ARGS = NodeKind("ND_FUNCALL_ARGS") // value に引数たちの abstSyntaxNode のスライス
	ND_FUNDEF       = NodeKind("ND_FUNDEF")       // value に関数名、leftNode に ND_FUNDEF_ARGS, rightNode に 関数のstmt
	ND_FUNDEF_ARGS  = NodeKind("ND_FUNDEF_ARGS")  // value に args の node のスライスを詰める
)

type TypeKind string

const (
	TypeInt = TypeKind("int")
	TypePtr = TypeKind("pointer")

	PointerSize = 8
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
	nodeKind  NodeKind
	leftNode  *abstSyntaxNode
	rightNode *abstSyntaxNode
	value     interface{} // num の値や、local variable の offset を入れる
}

func makeNewAbstSyntaxNode(nodeKind NodeKind, leftNode *abstSyntaxNode, rightNode *abstSyntaxNode, value interface{}) *abstSyntaxNode {
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

func getSizeOf(node *abstSyntaxNode) (size int, err error) {
	switch node.nodeKind {
	case ND_NUM:
		return 4, nil
	case ND_LVAR:
		nvar := node.value.(variable)
		switch nvar.kind {
		case TypeInt:
			return 4, nil
		case TypePtr:
			return 8, nil
		default:
			return 0, fmt.Errorf("getSizeOf: %s is not implemented", string(nvar.kind))
		}
	default:
		return 0, fmt.Errorf("getSizeOf: %s is not a valid type", string(node.nodeKind))
	}

	return 0, nil
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
