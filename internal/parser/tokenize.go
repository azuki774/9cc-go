package parser

import (
	"fmt"
	"strconv"
)

type TokenKind string

const (
	// Token.kind に入る
	TK_UNDEFINED = TokenKind("TK_UNDEFINED")
	TK_NUM       = TokenKind("TK_NUM") // Token.Value -> int
	TK_ADD       = TokenKind("TK_ADD")
	TK_SUB       = TokenKind("TK_SUB")
	TK_MUL       = TokenKind("TK_MUL")
	TK_DIV       = TokenKind("TK_DIV")
	TK_AND       = TokenKind("TK_AND")
	TK_LEFTPAT   = TokenKind("TK_LEFTPAT")
	TK_RIGHTPAT  = TokenKind("TK_RIGHTPAT")
	TK_BLOCKL    = TokenKind("TK_BLOCKL")
	TK_BLOCKR    = TokenKind("TK_BLOCKR")
	TK_COMP      = TokenKind("TK_COMP")      // ==
	TK_NOTEQ     = TokenKind("TK_NOTEQ")     // !=
	TK_LT        = TokenKind("TK_LT")        // <
	TK_LTQ       = TokenKind("TK_LTQ")       // <=
	TK_GT        = TokenKind("TK_GT")        // >
	TK_GTQ       = TokenKind("TK_GTQ")       // >=
	TK_EQ        = TokenKind("TK_EQ")        // =
	TK_SEMICOLON = TokenKind("TK_SEMICOLON") // ;
	TK_COMMA     = TokenKind("TK_COMMA")
	TK_RETURN    = TokenKind("TK_RETURN") // return
	TK_IF        = TokenKind("TK_IF")
	TK_ELSE      = TokenKind("TK_ELSE")
	TK_WHILE     = TokenKind("TK_WHILE")
	TK_FOR       = TokenKind("TK_FOR")
	TK_IDENT     = TokenKind("TK_IDENT") // Token.Value -> string (name)
	TK_SIZEOF    = TokenKind("TK_SIZEOF")
	TK_EOF       = TokenKind("TK_EOF")

	// // Token category in kind
	// TK_LIST_IDENT = []int{TK_IDENT}
)

var (
	// parser用
	TK_SYMBOL_LIST = []byte{BYTE_SYMBOL_ADD, BYTE_SYMBOL_SUB, BYTE_SYMBOL_MUL, BYTE_SYMBOL_DIV, BYTE_LEFTPAT, BYTE_RIGHTPAT, BYTE_EQUAL, BYTE_EXC, BYTE_LT, BYTE_GT, BYTE_SEMICOLON, BYTE_BLOCKL, BYTE_BLOCKR, BYTE_COMMA, BYTE_AND}
	TK_SPACE       = []byte{BYTE_SPACE}                             // スペース
	TK_DIGIT       = []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57} // 0 - 9
)

var (
	BYTE_SYMBOL_ADD = byte(43)
	BYTE_SYMBOL_SUB = byte(45)
	BYTE_SYMBOL_MUL = byte(42)
	BYTE_SYMBOL_DIV = byte(47)
	BYTE_LEFTPAT    = byte(40)
	BYTE_RIGHTPAT   = byte(41)
	BYTE_AND        = byte(38)
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
	BYTE_BLOCKL     = byte(123)
	BYTE_BLOCKR     = byte(125)
	BYTE_COMMA      = byte(44) // ,
)

type Token struct {
	kind  TokenKind
	value interface{}
}

func (token *Token) ShowString() (str string) {
	switch token.kind {
	case TK_UNDEFINED:
		str = "TK_UNDEFINED"
	case TK_NUM:
		str = fmt.Sprintf("TK_NUM: %d", token.value.(int))
	case TK_IDENT:
		str = fmt.Sprintf("TK_IDENT: %s", token.value.(string))
	case TK_EOF:
		str = "TK_EOF"
	default:
		str = string(token.kind)
	}
	return str
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
	token.kind = TK_UNDEFINED

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
			switch word {
			case "return":
				token = Token{kind: TK_RETURN}
			case "if":
				token = Token{kind: TK_IF}
			case "else":
				token = Token{kind: TK_ELSE}
			case "while":
				token = Token{kind: TK_WHILE}
			case "for":
				token = Token{kind: TK_FOR}
			case "sizeof":
				token = Token{kind: TK_SIZEOF}
			default:
				token = Token{kind: TK_IDENT, value: word}
			}

			return
		}

		ss.nextChar() // ポインタだけすすめる

		// 処理対象がSymbolの場合
		if contains(nChar, TK_SYMBOL_LIST) {
			switch nChar {
			case BYTE_SYMBOL_ADD:
				token = Token{kind: TK_ADD}
			case BYTE_SYMBOL_SUB:
				token = Token{kind: TK_SUB}
			case BYTE_SYMBOL_MUL:
				token = Token{kind: TK_MUL}
			case BYTE_SYMBOL_DIV:
				token = Token{kind: TK_DIV}
			case BYTE_LEFTPAT:
				token = Token{kind: TK_LEFTPAT}
			case BYTE_RIGHTPAT:
				token = Token{kind: TK_RIGHTPAT}
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
			case BYTE_BLOCKL:
				token = Token{kind: TK_BLOCKL}
			case BYTE_BLOCKR:
				token = Token{kind: TK_BLOCKR}
			case BYTE_COMMA:
				token = Token{kind: TK_COMMA}
			case BYTE_AND:
				token = Token{kind: TK_AND}
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
