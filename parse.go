package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Parse a package directory into an AST and FileSet.
func Parse(path string) (*ast.Package, *token.FileSet, error) {
	var fs = token.NewFileSet()
	var pkgs, e = parser.ParseDir(fs, path, nil, 0)
	if e != nil {
		return nil, nil, e
	}
	for _, pkg := range pkgs {
		return pkg, fs, nil
	}
	return nil, nil, fmt.Errorf("missing package at %s", path)
}
