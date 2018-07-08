package main

import (
	"fmt"
	"go/ast"
	"strings"
)

// RewriteRules converts user input of T=T1 or T=*T1 into
// a map of *ast.Ident values and a replacement ast.Node for
// anything that matches the ident.
func RewriteRules(inputs []string) (map[string]ast.Node, error) {
	var results = make(map[string]ast.Node)
	for _, input := range inputs {
		var tuple = strings.Split(input, "=")
		if len(tuple) != 2 {
			return nil, fmt.Errorf("invalid rewrite rule: %s", input)
		}
		var src = tuple[0]
		var dest = tuple[1]
		if dest[0] == '*' {
			results[src] = &ast.StarExpr{
				X: &ast.Ident{
					Name: dest[1:],
				},
			}
			continue
		}
		results[src] = &ast.Ident{
			Name: dest,
		}
	}
	return results, nil
}
