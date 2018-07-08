# gogen #
**AST rewrite powered generics for Go**

## Usage ##

Somewhere in your code you need to define a generic type. Something like:

```golang
type Generic interface{}
```

From here, you can write any code that uses the generic type:

```golang
func Add(left, right Generic) Generic {
    return left + right
}
```

Using the `gogen` CLI you can then render the generic code into a concrete type:

```bash
gogen -source github.com/username/package -destination . -rewrite "Generic=int"
```

The source code will be rewritten without the `type Generic interface{}` line and with all references to `Generic` will be replaced with the concrete type.

## More Options

In addition to direct rewrites, the CLI also offers the following options:

*   -unstar T
    
    This option triggers the rewrite of `*T` within the source to `T` before being rewritten as the concrete type.

*   -rewrite "T=*int"

    If the concrete type given in a rewrite rule is prefixed with an `*` then a pointer type will be used.

*   -package newname

    The package option causes the renderer to set a custom package name in all files rendered. By default, the original package name is reused.