package parser

import (
	"reflect"
	"testing"
)

func Test_contains(t *testing.T) {
	type args struct {
		b  byte
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{b: 1, bs: []byte{1, 2, 3}},
			want: true,
		},
		{
			name: "test2",
			args: args{b: 1, bs: []byte{11, 22, 33}},
			want: false,
		},
		{
			name: "plus",
			args: args{b: byte('+'), bs: []byte{1, 2, 43}},
			want: true,
		},
		{
			name: "minus",
			args: args{b: byte('-'), bs: []byte{1, 2, 45}},
			want: true,
		},
		{
			name: "multiple",
			args: args{b: byte('*'), bs: []byte{1, 2, 42}},
			want: true,
		},
		{
			name: "div",
			args: args{b: byte('/'), bs: []byte{1, 2, 47}},
			want: true,
		},
		{
			name: "space",
			args: args{b: byte(' '), bs: []byte{1, 2, 32}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.b, tt.args.bs); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextToken(t *testing.T) {
	ss1 := newStringStream("     123")
	ss2 := newStringStream("+")
	ss3 := newStringStream("  - ")
	ss4 := newStringStream("  * ")
	ss5 := newStringStream("/ ")
	ss6 := newStringStream("()")
	ss7 := newStringStream(")")
	ss8 := newStringStream("== ")
	ss9 := newStringStream("!=")
	ss10 := newStringStream("< 3")
	ss11 := newStringStream("<= 2")
	ss12 := newStringStream(">1")
	ss13 := newStringStream(">=10")
	ss14 := newStringStream("=10")
	ss15 := newStringStream(";")
	ss16 := newStringStream("a=123")
	ss17 := newStringStream("ABC=123")
	ss18 := newStringStream("return 1;")
	ss19 := newStringStream("if a")
	ss20 := newStringStream("else")
	ss21 := newStringStream("while")
	ss22 := newStringStream("for")
	ss23 := newStringStream(";;")
	ss24 := newStringStream("{xxxxx")
	ss25 := newStringStream("}yyyyy")
	ss26 := newStringStream(",abcd")
	ss27 := newStringStream("&abcd")
	ss28 := newStringStream("*x")
	type args struct {
		ss *stringStream
	}
	tests := []struct {
		name      string
		args      args
		wantToken Token
		wantErr   bool
	}{
		{
			name:      "num",
			args:      args{ss: ss1},
			wantToken: Token{kind: TK_NUM, value: 123},
			wantErr:   false,
		},
		{
			name:      "plus",
			args:      args{ss: ss2},
			wantToken: Token{kind: TK_SYMBOL_ADD},
			wantErr:   false,
		},
		{
			name:      "minus",
			args:      args{ss: ss3},
			wantToken: Token{kind: TK_SYMBOL_SUB},
			wantErr:   false,
		},
		{
			name:      "mul",
			args:      args{ss: ss4},
			wantToken: Token{kind: TK_SYMBOL_MUL},
			wantErr:   false,
		},
		{
			name:      "div",
			args:      args{ss: ss5},
			wantToken: Token{kind: TK_SYMBOL_DIV},
			wantErr:   false,
		},
		{
			name:      "left (",
			args:      args{ss: ss6},
			wantToken: Token{kind: TK_SYMBOL_LEFTPAT},
			wantErr:   false,
		},
		{
			name:      "right )",
			args:      args{ss: ss7},
			wantToken: Token{kind: TK_SYMBOL_RIGHTPAT},
			wantErr:   false,
		},
		{
			name:      "==",
			args:      args{ss: ss8},
			wantToken: Token{kind: TK_COMP},
			wantErr:   false,
		},
		{
			name:      "!=",
			args:      args{ss: ss9},
			wantToken: Token{kind: TK_NOTEQ},
			wantErr:   false,
		},
		{
			name:      "<",
			args:      args{ss: ss10},
			wantToken: Token{kind: TK_LT},
			wantErr:   false,
		},
		{
			name:      "<=",
			args:      args{ss: ss11},
			wantToken: Token{kind: TK_LTQ},
			wantErr:   false,
		},
		{
			name:      ">",
			args:      args{ss: ss12},
			wantToken: Token{kind: TK_GT},
			wantErr:   false,
		},
		{
			name:      ">=",
			args:      args{ss: ss13},
			wantToken: Token{kind: TK_GTQ},
			wantErr:   false,
		},
		{
			name:      "=",
			args:      args{ss: ss14},
			wantToken: Token{kind: TK_EQ},
			wantErr:   false,
		},
		{
			name:      ";",
			args:      args{ss: ss15},
			wantToken: Token{kind: TK_SEMICOLON},
			wantErr:   false,
		},
		{
			name:      "a",
			args:      args{ss: ss16},
			wantToken: Token{kind: TK_IDENT, value: "a"},
			wantErr:   false,
		},
		{
			name:      "ABC",
			args:      args{ss: ss17},
			wantToken: Token{kind: TK_IDENT, value: "ABC"},
			wantErr:   false,
		},
		{
			name:      "return",
			args:      args{ss: ss18},
			wantToken: Token{kind: TK_RETURN},
			wantErr:   false,
		},
		{
			name:      "if",
			args:      args{ss: ss19},
			wantToken: Token{kind: TK_IF},
			wantErr:   false,
		},
		{
			name:      "else",
			args:      args{ss: ss20},
			wantToken: Token{kind: TK_ELSE},
			wantErr:   false,
		},
		{
			name:      "while",
			args:      args{ss: ss21},
			wantToken: Token{kind: TK_WHILE},
			wantErr:   false,
		},
		{
			name:      "for",
			args:      args{ss: ss22},
			wantToken: Token{kind: TK_FOR},
			wantErr:   false,
		},
		{
			name:      ";;",
			args:      args{ss: ss23},
			wantToken: Token{kind: TK_SEMICOLON},
			wantErr:   false,
		},
		{
			name:      "{",
			args:      args{ss: ss24},
			wantToken: Token{kind: TK_BLOCKL},
			wantErr:   false,
		},
		{
			name:      "}",
			args:      args{ss: ss25},
			wantToken: Token{kind: TK_BLOCKR},
			wantErr:   false,
		},
		{
			name:      ",",
			args:      args{ss: ss26},
			wantToken: Token{kind: TK_COMMA},
			wantErr:   false,
		},
		{
			name:      "&",
			args:      args{ss: ss27},
			wantToken: Token{kind: TK_SYMBOL_AND},
			wantErr:   false,
		},
		{
			name:      "*x",
			args:      args{ss: ss28},
			wantToken: Token{kind: TK_SYMBOL_MUL},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := getNextToken(tt.args.ss)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNextToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("getNextToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
