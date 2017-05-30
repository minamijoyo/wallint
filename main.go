package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
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

type file struct {
	astFile  *ast.File
	fset     *token.FileSet
	src      []byte
	filename string
}

type problem struct {
	position token.Position
	text     string
}

func lint(filename string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to parse file: %s", filename)
	}
	// pp.Println(astFile)

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to read file: %s", filename)
	}

	f := &file{
		astFile:  astFile,
		fset:     fset,
		src:      src,
		filename: filename,
	}

	f.lint()
	return nil
}

func (f *file) lint() {

	check := func(ident *ast.Ident) {
		name := ident.Name
		if keywords[name] {
			text := fmt.Sprintf("%s was found\n", name)
			p := f.errorf(ident, text)
			fmt.Printf("%v: %s", p.position, p.text)
		}
	}

	visitor := func(node ast.Node) bool {
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

	ast.Walk(walker(visitor), f.astFile)
}

func (f *file) errorf(n ast.Node, text string) *problem {
	pos := f.fset.Position(n.Pos())
	problem := &problem{
		position: pos,
		text:     text,
	}
	return problem
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

var keywords = map[string]bool{
	"wally": true,
}
