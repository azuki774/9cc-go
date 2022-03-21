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
	ss1 := newStream("     123")
	ss2 := newStream("+")
	ss3 := newStream("  - ")

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
			wantToken: Token{kind: TK_SYMBOL, value: TK_SYMBOL_LIST[0]},
			wantErr:   false,
		},
		{
			name:      "minus",
			args:      args{ss: ss3},
			wantToken: Token{kind: TK_SYMBOL, value: TK_SYMBOL_LIST[1]},
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
