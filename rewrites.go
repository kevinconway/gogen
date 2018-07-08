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
		var star = dest[0] == '*'
		if star {
			dest = dest[1:]
		}
		var id ast.Node = &ast.Ident{
			Name: dest,
		}
		if strings.Contains(dest, ".") {
			var parts = strings.Split(dest, ".")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid rewrite rule: %s", input)
			}
			id = &ast.SelectorExpr{
				Sel: &ast.Ident{Name: parts[1]},
				X:   &ast.Ident{Name: parts[0]},
			}
		}
		if star {
			id = &ast.StarExpr{
				X: id.(ast.Expr),
			}
		}
		results[src] = id
	}
	return results, nil
}
