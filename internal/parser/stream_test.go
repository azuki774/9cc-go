package parser

import (
	"reflect"
	"testing"
)

func Test_stringStream_ok(t *testing.T) {
	type fields struct {
		str   string
		index int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test1",
			fields: fields{str: "0123456789", index: 0},
			want:   true,
		},
		{
			name:   "test2",
			fields: fields{str: "0123456789", index: 9},
			want:   true,
		},
		{
			name:   "test3",
			fields: fields{str: "0123456789", index: 10},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStream{
				str:   tt.fields.str,
				index: tt.fields.index,
			}
			if got := ss.ok(); got != tt.want {
				t.Errorf("stringStream.ok() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringStream_nextChar(t *testing.T) {
	type fields struct {
		str   string
		index int
	}
	tests := []struct {
		name   string
		fields fields
		wantB  byte
	}{
		{
			name:   "test1",
			fields: fields{str: "0123456789", index: 0},
			wantB:  byte('0'),
		},
		{
			name:   "test2",
			fields: fields{str: "0123456789", index: 8},
			wantB:  byte('8'),
		},
		{
			name:   "test3",
			fields: fields{str: "0123456789", index: 9},
			wantB:  byte('9'),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStream{
				str:   tt.fields.str,
				index: tt.fields.index,
			}
			if gotB := ss.nextChar(); gotB != tt.wantB {
				t.Errorf("stringStream.nextChar() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func Test_stringStream_nextPeekChar(t *testing.T) {
	type fields struct {
		str   string
		index int
	}
	tests := []struct {
		name   string
		fields fields
		wantB  byte
	}{
		{
			name:   "test1",
			fields: fields{str: "0123456789", index: 0},
			wantB:  byte('0'),
		},
		{
			name:   "test2",
			fields: fields{str: "0123456789", index: 8},
			wantB:  byte('8'),
		},
		{
			name:   "test3",
			fields: fields{str: "0123456789", index: 9},
			wantB:  byte('9'),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStream{
				str:   tt.fields.str,
				index: tt.fields.index,
			}
			if gotB := ss.nextPeekChar(); gotB != tt.wantB {
				t.Errorf("stringStream.nextPeekChar() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func Test_tokenStream_ok(t *testing.T) {
	type fields struct {
		tokens []Token
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test1",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 0},
			want:   true,
		},
		{
			name:   "test2",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 1},
			want:   true,
		},
		{
			name:   "test3",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 3},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &tokenStream{
				tokens: tt.fields.tokens,
				index:  tt.fields.index,
			}
			if got := ts.ok(); got != tt.want {
				t.Errorf("tokenStream.ok() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenStream_nextToken(t *testing.T) {
	type fields struct {
		tokens []Token
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		wantTk Token
	}{
		{
			name:   "test1",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 0},
			wantTk: Token{kind: TK_SYMBOL_LEFTPAT},
		},
		{
			name:   "test2",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 1},
			wantTk: Token{kind: TK_NUM, value: 123},
		},
		{
			name:   "test3",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 3},
			wantTk: Token{kind: TK_EOF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &tokenStream{
				tokens: tt.fields.tokens,
				index:  tt.fields.index,
			}
			if gotTk := ts.nextToken(); !reflect.DeepEqual(gotTk, tt.wantTk) {
				t.Errorf("tokenStream.nextToken() = %v, want %v", gotTk, tt.wantTk)
			}
		})
	}
}

func Test_tokenStream_nextPeekToken(t *testing.T) {
	type fields struct {
		tokens []Token
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		wantTk Token
	}{
		{
			name:   "test1",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 0},
			wantTk: Token{kind: TK_SYMBOL_LEFTPAT},
		},
		{
			name:   "test2",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 1},
			wantTk: Token{kind: TK_NUM, value: 123},
		},
		{
			name:   "test3",
			fields: fields{tokens: []Token{{kind: TK_SYMBOL_LEFTPAT}, {kind: TK_NUM, value: 123}, {kind: TK_SYMBOL_RIGHTPAT}, {kind: TK_EOF}}, index: 3},
			wantTk: Token{kind: TK_EOF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &tokenStream{
				tokens: tt.fields.tokens,
				index:  tt.fields.index,
			}
			if gotTk := ts.nextPeekToken(); !reflect.DeepEqual(gotTk, tt.wantTk) {
				t.Errorf("tokenStream.nextPeekToken() = %v, want %v", gotTk, tt.wantTk)
			}
		})
	}
}

func Test_stringStream_isNextWordCharNoPrefix(t *testing.T) {
	ss1 := newStringStream("hoge")
	ss2 := newStringStream("HuGa")
	ss3 := newStringStream("0123")
	ss4 := newStringStream("_")
	ss5 := newStringStream(";")
	ss6 := newStringStream("(")
	ss7 := newStringStream("-")
	tests := []struct {
		name string
		ss   *stringStream
		want bool
	}{
		{
			name: "test1",
			ss:   ss1,
			want: true,
		},
		{
			name: "test2",
			ss:   ss2,
			want: true,
		},
		{
			name: "test3",
			ss:   ss3,
			want: true,
		},
		{
			name: "test4",
			ss:   ss4,
			want: true,
		},
		{
			name: "test5",
			ss:   ss5,
			want: false,
		},
		{
			name: "test6",
			ss:   ss6,
			want: false,
		},
		{
			name: "test7",
			ss:   ss7,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.isNextWordCharNoPrefix(); got != tt.want {
				t.Errorf("stringStream.isNextWordCharNoPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringStream_nextWord(t *testing.T) {
	ss1 := newStringStream("hoge")
	ss2 := newStringStream("fuga0123")
	ss3 := newStringStream("return;")
	ss4 := newStringStream("abc+def")
	ss5 := newStringStream("sepa = rate")
	tests := []struct {
		name  string
		ss    *stringStream
		wantS string
	}{
		{
			name:  "test1",
			ss:    ss1,
			wantS: "hoge",
		},
		{
			name:  "test2",
			ss:    ss2,
			wantS: "fuga0123",
		},
		{
			name:  "test3",
			ss:    ss3,
			wantS: "return",
		},
		{
			name:  "test4",
			ss:    ss4,
			wantS: "abc",
		},
		{
			name:  "test5",
			ss:    ss5,
			wantS: "sepa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := tt.ss.nextWord(); gotS != tt.wantS {
				t.Errorf("stringStream.nextWord() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
