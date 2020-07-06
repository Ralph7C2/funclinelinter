package linter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/fatih/color"
)

type fixer struct {
	fileSet   *token.FileSet
	node      ast.Node
	printedTo token.Pos
}

func Fix(filename string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.DeclarationErrors)
	if err != nil {
		panic(fmt.Sprintf("Parser error: %s", err))
	}
	fx := &fixer{fileSet: fset, node: f}
	fmt.Print("package ")
	ast.Walk(fx, f)
}

type fnPrinter struct {
	fileSet   *token.FileSet
	printedTo token.Pos
}

func (p *fnPrinter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return p
	}
	if node.Pos() < p.printedTo {
		return p
	}
	buf := bytes.NewBuffer([]byte{})
	printer.Fprint(buf, p.fileSet, node)
	fmt.Println(buf.String())
	p.printedTo = node.End()
	return p
}

func (f *fixer) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return f
	}
	if node.Pos() < f.printedTo {
		return f
	}
	fcc := &funcChildChecker{}
	ast.Walk(fcc, node)
	if fcc.foundFunc {
		color.Set(color.FgRed)
		switch fn := node.(type) {
		case *ast.FuncLit:
			ast.Walk(&fnPrinter{fileSet: f.fileSet}, fn)
			f.printedTo = fn.End()
		case *ast.FuncType:
			ast.Walk(&fnPrinter{fileSet: f.fileSet}, fn)
			f.printedTo = fn.End()
		case *ast.FuncDecl:
			ast.Walk(&fnPrinter{fileSet: f.fileSet}, fn)
			f.printedTo = fn.End()
		}
		color.Unset()
		return f
	}
	printer.Fprint(os.Stdout, f.fileSet, node)
	fmt.Println("\n")
	f.printedTo = node.End()
	return f
}

type funcChildChecker struct {
	foundFunc bool
}

func (f *funcChildChecker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return f
	}
	switch node.(type) {
	case *ast.FuncDecl:
		f.foundFunc = true
		return nil
	case *ast.FuncType:
		f.foundFunc = true
		return nil
	case *ast.FuncLit:
		f.foundFunc = true
		return nil
	}
	return f
}
