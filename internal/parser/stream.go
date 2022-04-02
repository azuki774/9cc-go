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

// 今の文字が2文字目以降の変数名、キーワードとして使えるかどうか
func (ss *stringStream) isNextWordCharNoPrefix() bool {
	if !ss.ok() {
		return false
	}

	b := ss.str[ss.index]

	if BYTE_a <= b && b <= BYTE_z { // a <= var b <= z
		return true
	}

	if BYTE_A <= b && b <= BYTE_Z { // A <= var b <= Z
		return true
	}

	if contains(b, TK_DIGIT) {
		return true
	}

	if b == BYTE_UNDERBAR {
		return true
	}

	return false
}

// この位置から次の変数、キーワードとして使えるWordを取り出す
func (ss *stringStream) nextWord() (s string) {
	if !ss.ok() {
		return ""
	}

	// 1文字目
	b := ss.str[ss.index]

	if (BYTE_a <= b && b <= BYTE_z) || (BYTE_A <= b && b <= BYTE_Z) { // a <= var b <= z or A <= var b <= Z
		s += string(b)
		ss.nextChar()
	} else {
		// 1文字が既に無効
		return ""
	}

	for {
		if ss.isNextWordCharNoPrefix() {
			// 2文字目以降が続きとして有効
			s += string(ss.nextChar())
		} else {
			break
		}
	}

	return s
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
