package parser

type stringStream struct {
	str   string
	index int
}

type tokenStream struct {
	tokens []Token
	index  int
}

func (ss *stringStream) ok() bool {
	return ss.index+1 <= len(ss.str)
}

func (ss *stringStream) nextChar() (b byte) {
	b = ss.str[ss.index]
	ss.index++
	return b
}

func (ss *stringStream) nextPeekChar() (b byte) {
	b = ss.str[ss.index]
	return b
}

func (ts *tokenStream) ok() bool {
	return !(ts.tokens[ts.index].kind == TK_EOF)
}

func (ts *tokenStream) nextToken() (tk Token) {
	tk = ts.tokens[ts.index]
	ts.index++
	return tk
}

func (ts *tokenStream) nextPeekToken() (tk Token) {
	tk = ts.tokens[ts.index]
	return tk
}

func newTokenStream(tokens []Token) *tokenStream {
	return &tokenStream{tokens: tokens, index: 0}
}

func newStringStream(str string) *stringStream {
	return &stringStream{str: str, index: 0}
}
