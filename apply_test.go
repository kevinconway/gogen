package main

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

var source = `package main

import (
	"testing"
)

type T interface{}
type T2 interface{}
type T3 interface{}

type S struct {
	T  T
	T2 T2
	T3 []T3
}

func TestSomething(t *testing.T) {
	t.Skip("fake test")
}

func DoSomething(t T, tt T2) (T, T2, error) {
	var a T
	var b T2
	return a, b, nil
}

func DoSomethingStar(t *T, tt *T2) (*T, *T2, error) {
	var a *T
	var b *T2
	return a, b, nil
}
`

var result = `package main

import (
	"testing"
)

type S struct {
	T  int
	T2 string
	T3 []bool
}

func TestSomething(t *testing.T) {
	t.Skip("fake test")
}

func DoSomething(t int, tt string) (int, string, error) {
	var a int
	var b string
	return a, b, nil
}

func DoSomethingStar(t *int, tt string) (*int, string, error) {
	var a *int
	var b string
	return a, b, nil
}
`

func TestApply(t *testing.T) {
	var fs = token.NewFileSet()
	var n, e = parser.ParseFile(fs, "test.go", source, 0)
	if e != nil {
		t.Fatal(e.Error())
	}
	var rewrites, _ = RewriteRules([]string{"T=int", "T2=string", "T3=bool"})
	var unstars = []string{"T2"}
	var apply = ApplyFN(rewrites, unstars)
	var final = astutil.Apply(n, apply, nil)
	var resultBuffer = bytes.NewBufferString("")
	_ = format.Node(resultBuffer, fs, final)
	if resultBuffer.String() != result {
		t.Errorf("Expected:\n%s\n\nReceived:\n%s", result, resultBuffer.String())
	}
}
