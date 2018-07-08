package main

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestRewriteRules(t *testing.T) {
	type args struct {
		inputs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]ast.Node
		wantErr bool
	}{
		{
			name:    "invalid",
			args:    args{inputs: []string{"noequal"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nostar",
			args: args{inputs: []string{"T=NT"}},
			want: map[string]ast.Node{
				"T": &ast.Ident{Name: "NT"},
			},
			wantErr: false,
		},
		{
			name: "star",
			args: args{inputs: []string{"T=*NT"}},
			want: map[string]ast.Node{
				"T": &ast.StarExpr{X: &ast.Ident{Name: "NT"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RewriteRules(tt.args.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RewriteRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RewriteRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
