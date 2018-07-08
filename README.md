# gogen #
**AST rewrite powered generics for Go**

## Usage ##

Somewhere in your code you need to define one or more generic types. Something like:

```golang
type I interface{}
type F interface{}
```

From here, you can write any code that uses those generic type (including tests):

```golang
func Add(left, right I) I {
    return left + right
}
func Sub(left, right F) F {
    return left - right
}
```

Using the `gogen` CLI you can then render the generic code with concrete types:

```bash
gogen -source github.com/username/package -destination . -rewrite "I=uint16" -rewrite "F=float64
```

The source code will be rewritten without the `type I interface{}` lines and
with all references to generic types replaced with concrete types. The output
will be formatted with `gofmt` (but not `goimports`).

## More Options

In addition to direct rewrites, the CLI also offers the following options:

*   -unstar T
    
    This option triggers the rewrite of `*T` within the source to `T` before
    being rewritten as the concrete type. This is useful if rendering a generic
    as an interface type (such as `-rewrite "T=io.Reader"). The option may be
    given multiple times to unstar multiple generic types.

*   -rewrite "T=*int"

    If the concrete type given in a rewrite rule is prefixed with an `*` then
    a pointer type will be used. If the generic type is already a pointer then
    this could result in a double pointer to the concrete type. Use the
    `unstar` flag to prevent this. Any number of rewrites may be given per
    invocation.

*   -rewrite "T=pkg.Name"

    If the concrete type has a package prefix then the resulting code will
    correctly select the type from the target package. However, the rewrite
    will not change import statements or add new ones. It is recommended to
    use something like `goimports` to correct any missing or excess import
    statements after rewrites.

*   -package newname

    The package option causes the renderer to set a custom package name in
    all files rendered. By default, the original package name is reused.