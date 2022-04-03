package parser

import (
	"fmt"
	"strconv"
)

var (
	// Token.kind に入る
	TK_UNDEFINED       = 0
	TK_NUM             = 1 // Token.Value -> int
	TK_SYMBOL_ADD      = 11
	TK_SYMBOL_SUB      = 12
	TK_SYMBOL_MUL      = 13
	TK_SYMBOL_DIV      = 14
	TK_SYMBOL_LEFTPAT  = 15
	TK_SYMBOL_RIGHTPAT = 16
	TK_COMP            = 21  // ==
	TK_NOTEQ           = 22  // !=
	TK_LT              = 23  // <
	TK_LTQ             = 24  // <=
	TK_GT              = 25  // >
	TK_GTQ             = 26  // >=
	TK_EQ              = 27  // =
	TK_SEMICOLON       = 31  // ;
	TK_IDENT           = 101 // Token.Value -> string (name)
	TK_EOF             = 255

	// Token category in kind
	TK_LIST_IDENT = []int{TK_IDENT}
)

var (
	// parser用
	TK_SYMBOL_LIST = []byte{BYTE_SYMBOL_ADD, BYTE_SYMBOL_SUB, BYTE_SYMBOL_MUL, BYTE_SYMBOL_DIV, BYTE_LEFTPAT, BYTE_RIGHTPAT, BYTE_EQUAL, BYTE_EXC, BYTE_LT, BYTE_GT, BYTE_SEMICOLON}
	TK_SPACE       = []byte{BYTE_SPACE}                             // スペース
	TK_DIGIT       = []byte{47, 48, 49, 50, 51, 52, 53, 54, 55, 56} // 0 - 9
)

var (
	BYTE_SYMBOL_ADD = byte(43)
	BYTE_SYMBOL_SUB = byte(45)
	BYTE_SYMBOL_MUL = byte(42)
	BYTE_SYMBOL_DIV = byte(47)
	BYTE_LEFTPAT    = byte(40)
	BYTE_RIGHTPAT   = byte(41)
	BYTE_SPACE      = byte(32)
	BYTE_EQUAL      = byte(61)
	BYTE_EXC        = byte(33)
	BYTE_LT         = byte(60) // <
	BYTE_GT         = byte(62) // >
	BYTE_SEMICOLON  = byte(59) // ;
	BYTE_UNDERBAR   = byte(95) // _
	BYTE_a          = byte(97)
	BYTE_A          = byte(65)
	BYTE_z          = byte(122)
	BYTE_Z          = byte(90)
)

type Token struct {
	kind  int
	value interface{}
}

func (token *Token) Show() {
	switch token.kind {
	case TK_UNDEFINED:
		fmt.Printf("TK_UNDEFINED\n")
	case TK_NUM:
		fmt.Printf("TK_NUM : %d\n", token.value.(int))
	case TK_IDENT:
		fmt.Printf("TK_IDENT : %s\n", token.value.(string))
	case TK_EOF:
		fmt.Printf("TK_EOF\n")
	default:
		fmt.Printf("TK_SYMBOL : %d\n", token.kind)
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
		nChar := ss.nextPeekChar()
		if !isContinueLoadNextChar(nChar, token) {
			break
		}

		if (BYTE_a <= nChar && nChar <= BYTE_z) || (BYTE_A <= nChar && nChar <= BYTE_Z) { // a <= nChar <= z or A <= nChar <= Z
			word := ss.nextWord()
			token = Token{kind: TK_IDENT, value: word}
			break
		}

		ss.nextChar() // ポインタだけすすめる

		// 処理対象がSymbolの場合
		if contains(nChar, TK_SYMBOL_LIST) {
			switch nChar {
			case BYTE_SYMBOL_ADD:
				token = Token{kind: TK_SYMBOL_ADD}
			case BYTE_SYMBOL_SUB:
				token = Token{kind: TK_SYMBOL_SUB}
			case BYTE_SYMBOL_MUL:
				token = Token{kind: TK_SYMBOL_MUL}
			case BYTE_SYMBOL_DIV:
				token = Token{kind: TK_SYMBOL_DIV}
			case BYTE_LEFTPAT:
				token = Token{kind: TK_SYMBOL_LEFTPAT}
			case BYTE_RIGHTPAT:
				token = Token{kind: TK_SYMBOL_RIGHTPAT}
			case BYTE_EQUAL:
				nnChar := ss.nextPeekChar()
				if nnChar == BYTE_EQUAL {
					ss.nextChar() // =
					token = Token{kind: TK_COMP}
				} else {
					token = Token{kind: TK_EQ}
				}
			case BYTE_EXC:
				// !=
				nnChar := ss.nextPeekChar()
				if nnChar == BYTE_EQUAL {
					ss.nextChar() // =
					token = Token{kind: TK_NOTEQ}
				} else {
					return Token{}, fmt.Errorf("getNextToken : != tokenize error")
				}
			case BYTE_LT:
				nnChar := ss.nextPeekChar()
				if nnChar == BYTE_EQUAL {
					ss.nextChar()               // =
					token = Token{kind: TK_LTQ} // <=
				} else {
					token = Token{kind: TK_LT} // <
				}
			case BYTE_GT:
				nnChar := ss.nextPeekChar()
				if nnChar == BYTE_EQUAL {
					ss.nextChar()               // =
					token = Token{kind: TK_GTQ} // >=
				} else {
					token = Token{kind: TK_GT} // >
				}
			case BYTE_SEMICOLON:
				token = Token{kind: TK_SEMICOLON}
			}

			break
		}

		if contains(nChar, TK_DIGIT) { // 0 - 9
			if token.kind == TK_UNDEFINED || token.kind == TK_NUM {
				token.kind = TK_NUM
				numString += string(nChar)
			} else {
				return Token{}, fmt.Errorf("getNextToken : digit tokenize error")
			}
		}

	}

	// 読み込み後の後処理

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

// 今読もうとしている文字のct = TK_***が、今読もうとしているTokenの続きならtrue、そうでないならfalse
func isContinueLoadNextChar(b byte, token Token) bool {
	if token.kind == TK_UNDEFINED {
		return true
	}

	if b == BYTE_SPACE {
		// 既になんらかの文字を読み込んでいたらTokenの切れ目
		return token.kind == TK_UNDEFINED
	}

	if token.kind == TK_NUM && contains(b, TK_SYMBOL_LIST) {
		// ex. 1+
		return false
	}

	if token.kind != TK_NUM && !contains(b, TK_DIGIT) {
		// ex. +1
		return false
	}

	return true
}

func TokenizeMain(str string) (tokens []Token, err error) {
	ss := newStringStream(str)
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
