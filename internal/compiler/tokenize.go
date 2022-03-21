package compiler

import (
	"fmt"
	"strconv"
)

var (
	TK_UNDEFINED = 0
	TK_SYMBOL    = 1
	TK_NUM       = 2
	TK_EOF       = 255
)

var (
	TK_SYMBOL_LIST = []byte{43, 45}                                 // + , -
	TK_SPACE       = []byte{32}                                     // スペース
	TK_DIGIT       = []byte{47, 48, 49, 50, 51, 52, 53, 54, 55, 56} // 0 - 9
)

type Token struct {
	kind  int
	value interface{}
}

func (token *Token) Show() {
	switch token.kind {
	case TK_UNDEFINED:
		fmt.Printf("TK_UNDEFINED\n")
	case TK_SYMBOL:
		fmt.Printf("TK_SYMBOL : %s\n", string(token.value.(byte)))
	case TK_NUM:
		fmt.Printf("TK_NUM : %d\n", token.value.(int))
	}
}

// b が bs に含まれるかどうか []byte版
func contains(b byte, bs []byte) bool {
	for _, v := range bs {
		if v == b {
			return true
		}
	}
	return false
}

func getNextToken(ss *stringStream) (token Token, err error) {
	numString := ""

	for {
		if !ss.ok() {
			// これ以上文字を読み込めないとき
			if token.kind == TK_UNDEFINED {
				token = Token{kind: TK_EOF}
			}
			break
		}

		// 次の文字が見て読むべきかどうか判定
		nChar, err := ss.nextPeekChar()
		if err != nil {
			return Token{}, err
		}

		tokenCategory := getTokenCategory(nChar)

		if !isContinueLoadNextChar(tokenCategory, token) {
			break
		}

		ss.nextChar() // ポインタだけすすめる

		if tokenCategory == TK_SYMBOL {
			token = Token{kind: TK_SYMBOL, value: nChar}
			break
		}

		if tokenCategory == TK_NUM {
			if token.kind == TK_UNDEFINED || token.kind == TK_NUM {
				token.kind = TK_NUM
				numString += string(nChar)
			} else {
				return Token{}, fmt.Errorf("getNextToken : digit tokenize error")
			}
		}
	}

	// 数値のときはToken.Valueに数値を移す
	if token.kind == TK_NUM {
		num, err := strconv.Atoi(numString)
		if err != nil {
			return Token{}, fmt.Errorf("getNextToken : Atoi error : %w", err)
		}
		token.value = num
	}

	return token, nil
}

// このbyteがどの種類に属するかを判定して、TK_*** を返す
func getTokenCategory(b byte) (category int) {
	if contains(b, TK_SPACE) {
		return TK_UNDEFINED
	}
	if contains(b, TK_SYMBOL_LIST) {
		return TK_SYMBOL
	}
	if contains(b, TK_DIGIT) {
		return TK_NUM
	}

	return TK_UNDEFINED
}

// 今読もうとしている文字のct = TK_***が、今読もうとしているTokenの続きならtrue、そうでないならfalse
func isContinueLoadNextChar(ct int, token Token) bool {
	if token.kind == TK_UNDEFINED {
		return true
	}

	if ct != token.kind {
		return false
	}
	return true
}

func TokenizeMain(str string) (tokens []Token, err error) {
	ss := newStream(str)
	for {
		newToken, err := getNextToken(ss)
		if err != nil {
			return nil, fmt.Errorf("TokenizeMain : %w", err)
		}
		tokens = append(tokens, newToken)
		if newToken.kind == TK_EOF {
			break
		}

	}
	return tokens, nil
}
