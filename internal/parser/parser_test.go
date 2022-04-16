package parser

import (
	"reflect"
	"testing"
)

func Test_makeNewAbstSyntaxNode(t *testing.T) {
	type args struct {
		nodeKind  NodeKind
		leftNode  *abstSyntaxNode
		rightNode *abstSyntaxNode
		value     interface{}
	}
	tests := []struct {
		name string
		args args
		want *abstSyntaxNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeNewAbstSyntaxNode(tt.args.nodeKind, tt.args.leftNode, tt.args.rightNode, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeNewAbstSyntaxNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_variableManager_reset(t *testing.T) {
	tests := []struct {
		name string
		v    *variableManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.reset()
		})
	}
}

func Test_variableManager_add(t *testing.T) {
	v1 := variableManager{}
	v1.reset()

	type args struct {
		varname string
		kind    TypeKind
	}
	tests := []struct {
		name     string
		v        *variableManager
		args     args
		wantNvar variable
		wantErr  bool
	}{
		{
			name:     "test1",
			v:        &v1,
			args:     args{varname: "a", kind: TypeInt},
			wantNvar: variable{kind: TypeInt, ptrTo: nil, offset: 8},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNvar, err := tt.v.add(tt.args.varname, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("variableManager.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNvar, tt.wantNvar) {
				t.Errorf("variableManager.add() = %v, want %v", gotNvar, tt.wantNvar)
			}
		})
	}
}
