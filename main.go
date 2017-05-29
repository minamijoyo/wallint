package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	// "github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: minilint filename")
		os.Exit(1)
	}
	filename := os.Args[1]
	if err := lint(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func lint(filename string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to parse file: %s", filename)
	}
	// pp.Println(astFile)
	ast.Walk(walker(visitor), astFile)
	return nil
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func visitor(node ast.Node) bool {
	switch v := node.(type) {
	case *ast.AssignStmt:
		for _, expr := range v.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				check(ident)
			}
		}
	}
	return true
}

var keywords = map[string]bool{
	"wally": true,
}

func check(ident *ast.Ident) {
	name := ident.Name
	if keywords[name] {
		fmt.Println(name)
	}
}
