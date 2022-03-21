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
	TK_SYMBOL_LIST = []byte{43, 45} // + , -
	TK_SPACE       = []byte{32}     // スペース
	TK_DIGIT       = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

type Token struct {
	kind  int
	value interface{}
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
			break
		}

		nChar, err := ss.nextChar()
		if err != nil {
			return Token{}, err
		}

		if contains(nChar, TK_SPACE) {
			if token.kind == 0 {
				// まだ何も読み込んでいない場合
				continue
			} else {
				// 読み込んだあとの半角スペースの場合
				break
			}
		}

		if contains(nChar, TK_SYMBOL_LIST) {
			token = Token{kind: TK_SYMBOL, value: nChar}
			break
		}

		if contains(nChar, TK_DIGIT) {
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

func TokenizeMain(str string) (tokens []Token, err error) {
	return nil, nil
}
