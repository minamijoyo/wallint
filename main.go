package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/pkg/errors"
)

func main() {
	filename := os.Args[1]
	if err := lint(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type file struct {
	astFile  *ast.File
	fset     *token.FileSet
	filename string
}

func lint(filename string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to parse file: %s", filename)
	}

	f := &file{
		astFile:  astFile,
		fset:     fset,
		filename: filename,
	}

	f.lint()
	return nil
}

func (f *file) lint() {
	visitor := func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.AssignStmt:
			for _, expr := range v.Lhs {
				if ident, ok := expr.(*ast.Ident); ok {
					f.check(ident)
				}
			}
		}
		return true
	}

	ast.Walk(walker(visitor), f.astFile)
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func (f *file) check(ident *ast.Ident) {
	name := ident.Name
	if keywords[name] {
		pos := f.fset.Position(ident.Pos())
		text := fmt.Sprintf("%s was found\n", name)
		fmt.Printf("%v: %s", pos, text)
	}
}

var keywords = map[string]bool{
	"wally": true,
}
