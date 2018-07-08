package main

import (
	"go/ast"
	"go/format"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	var app = cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Value: ".",
			Usage: "import path of the package to render",
		},
		cli.StringSliceFlag{
			Name:  "rewrite, r",
			Usage: "a rewrite rule in the form of original=final or original=*final",
		},
		cli.StringSliceFlag{
			Name:  "unstar, u",
			Usage: "name of a source type for which to drop *",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: ".",
			Usage: "directory in which to write the results",
		},
		cli.StringFlag{
			Name:  "package, p",
			Usage: "package name to use for the resulting files",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		var pkg = filepath.Join(GOPATH(), "src", ctx.String("source"))
		var rewriteStrings = ctx.StringSlice("rewrite")
		var unstars = ctx.StringSlice("unstar")
		var destination, _ = filepath.Abs(ctx.String("destination"))
		var packageName = ctx.String("package")

		var p, fs, errPkg = Parse(pkg)
		if errPkg != nil {
			return cli.NewExitError(errPkg.Error(), 1)
		}
		var rewriteRules, rewriteErr = RewriteRules(rewriteStrings)
		if rewriteErr != nil {
			return cli.NewExitError(rewriteErr.Error(), 1)
		}
		p = astutil.Apply(p, ApplyFN(rewriteRules, unstars), nil).(*ast.Package)

		if packageName != "" {
			for _, f := range p.Files {
				f.Name.Name = packageName
			}
		}

		for srcPath, f := range p.Files {
			var outFile, errOutFile = os.Create(filepath.Join(destination, filepath.Base(srcPath)))
			if errOutFile != nil {
				return cli.NewExitError(errOutFile.Error(), 1)
			}
			_ = format.Node(outFile, fs, f)
			_ = outFile.Close()
		}
		return nil
	}
	var e = app.Run(os.Args)
	if e != nil {
		panic(e.Error())
	}
}
