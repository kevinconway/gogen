package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// ApplyGenDecl handles declaration statements. The purpose of this
// application is to remove generic declarations and their comments
func ApplyGenDecl(c *astutil.Cursor, gd *ast.GenDecl, rewrites map[string]ast.Node) bool {
	if len(gd.Specs) == 1 {
		if ts, tsOK := gd.Specs[0].(*ast.TypeSpec); tsOK {
			for name := range rewrites {
				if ts.Name.Name == name {
					c.Delete()
					return false
				}
			}
		}
	}
	return true
}

// ApplyField handles struct fields and function signature parameters.
// Internally, it switches on the star expression to implement the special
// case of "unstarring".
func ApplyField(c *astutil.Cursor, fd *ast.Field, rewrites map[string]ast.Node, unstars []string) bool {
	switch n := fd.Type.(type) {
	case *ast.StarExpr:
		return ApplyFieldStar(c, fd, n, rewrites, unstars)
	case *ast.Ident:
		return ApplyFieldIdent(c, fd, n, rewrites)
	default:
		return true
	}
}

// ApplyFieldStar handles all cases of pointers to T in structs and function
// signature definitions. If configured to unstar then it replaces the star
// expression with a simple ident node. Otherwise it replaces the ident node
// of the star expression.
func ApplyFieldStar(c *astutil.Cursor, fd *ast.Field, se *ast.StarExpr, rewrites map[string]ast.Node, unstars []string) bool {
	switch n := se.X.(type) {
	case *ast.Ident:
		for name, replacement := range rewrites {
			if n.Name == name {
				for _, unstar := range unstars {
					if n.Name == unstar {
						// Rewrite the star as a simple ident.
						fd.Type = replacement.(ast.Expr)
						c.Replace(fd)
						return true
					}
				}
				se.X = replacement.(ast.Expr)
				c.Replace(fd)
				return true
			}
		}
	default:
		return true
	}
	return true
}

// ApplyFieldIdent handles the case of any struct field or function parameter
// of type T by replacing it with the concrete type.
func ApplyFieldIdent(c *astutil.Cursor, fd *ast.Field, id *ast.Ident, rewrites map[string]ast.Node) bool {
	for name, replacement := range rewrites {
		if id.Name == name {
			fd.Type = replacement.(ast.Expr)
			c.Replace(fd)
			return true
		}
	}
	return true
}

// ApplyStar handles all general cases of a star expression. Pointers to the generic
// type are replaced.
func ApplyStar(c *astutil.Cursor, se *ast.StarExpr, rewrites map[string]ast.Node, unstars []string) bool {
	switch id := se.X.(type) {
	case *ast.Ident:
		for name, replacement := range rewrites {
			if id.Name == name {
				for _, unstar := range unstars {
					if id.Name == unstar {
						c.Replace(replacement)
						return true
					}
				}
				se.X = replacement.(ast.Expr)
				c.Replace(se)
				return true
			}
		}
	default:
		return true
	}
	return true
}

// ApplyIdent handles all generic cases of references a name.
func ApplyIdent(c *astutil.Cursor, id *ast.Ident, rewrites map[string]ast.Node) bool {
	switch n := c.Parent().(type) {
	case *ast.ValueSpec:
		// Ident nodes in a value spec that are not the type
		// are typically the variable name. Do not rewrite
		// these.
		if id != n.Type {
			return true
		}
	case *ast.SelectorExpr:
		// Ignore ident nodes that are child to a selector
		// expression to prevent accidental modification of
		// a remote type. For example, if the generic type is T
		// the code references otherpackage.T then we do not want
		// to rewrite the T into the concrete type because it isn't
		// the same.
		return true
	case *ast.Field:
		// Like VarSpec, ident nodes that are not the type of a field
		// are often the attribute name. Skip these.
		if id != n.Type {
			return true
		}
	}
	for name, replacement := range rewrites {
		if id.Name == name {
			fmt.Printf("%s %T\n", name, c.Parent())
			c.Replace(replacement)
			return true
		}
	}
	return true
}

// ApplyFN generates an Apply function for the AST rewrite.
func ApplyFN(rewrites map[string]ast.Node, unstars []string) func(*astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.GenDecl:
			return ApplyGenDecl(c, n, rewrites)
		case *ast.Field:
			return ApplyField(c, n, rewrites, unstars)
		case *ast.StarExpr:
			return ApplyStar(c, n, rewrites, unstars)
		case *ast.Ident:
			return ApplyIdent(c, n, rewrites)
		default:
			return true
		}
	}
}
