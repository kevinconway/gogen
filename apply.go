package main

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// ApplyFN generates an Apply function for the AST rewrite.
func ApplyFN(rewrites map[string]ast.Node, unstar []string) func(*astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		// First check if this line is the definition of type T generic.T
		// and remove it if it is.
		if gd, gdOK := c.Node().(*ast.GenDecl); gdOK {
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
		}
		// If not a definition then check if this is a param.
		if fd, fdOK := c.Node().(*ast.Field); fdOK {
			if se, seOK := fd.Type.(*ast.StarExpr); seOK {
				if id, idOK := se.X.(*ast.Ident); idOK {
					for name, replacement := range rewrites {
						if id.Name == name {
							for _, un := range unstar {
								if id.Name == un {
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
				}
			}
			if id, idOK := fd.Type.(*ast.Ident); idOK {
				for name, replacement := range rewrites {
					if id.Name == name {
						fd.Type = replacement.(ast.Expr)
						c.Replace(fd)
						return true
					}
				}
			}
		}
		// Next, identify the star expressions (*T) and apply
		// any rewrite rules there first.
		if se, seOK := c.Node().(*ast.StarExpr); seOK {
			if id, idOK := se.X.(*ast.Ident); idOK {
				for name, replacement := range rewrites {
					if id.Name == name {
						for _, un := range unstar {
							if id.Name == un {
								// Rewrite the star as a simple ident.
								c.Replace(replacement)
								return true
							}
						}
						se.X = replacement.(ast.Expr)
						c.Replace(se)
						return true
					}
				}
			}
		}
		// If not a star then check if it is a simple ident that needs
		// to be changed.
		if id, idOK := c.Node().(*ast.Ident); idOK {
			for name, replacement := range rewrites {
				if id.Name == name {
					switch c.Parent().(type) {
					case *ast.SelectorExpr:
					default:
						c.Replace(replacement)
						return true
					}
				}
			}
		}
		return true
	}
}
