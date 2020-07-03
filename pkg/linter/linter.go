package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

type linter struct {
	fset *token.FileSet
	out io.Writer
}

func NewLinter() *linter {
	return &linter{
		fset: token.NewFileSet(),
		out: os.Stdout,
	}
}

func (l *linter) Lint(fileName string) {
	f, err := parser.ParseFile(l.fset, fileName, nil, parser.DeclarationErrors)
	if err != nil {
		panic(fmt.Sprintf("Parser error: %s", err))
	}
	for s, object := range f.Scope.Objects {
		if object.Kind == ast.Fun {
			fd, ok := object.Decl.(*ast.FuncDecl)
			if !ok {
				fmt.Println("Error", s, "is ObjectKind Func, but decl is not FuncDecl?")
			}
			l.handleFunctionDefinition(fd)
		} else if object.Kind == ast.Typ {
			ts, ok := object.Decl.(*ast.TypeSpec)
			if !ok {
				fmt.Println("Error", s, "is ObjectKind Typ, but decl is not TypeSpec?")
			}
			l.handleTypeDefinition(ts)
		} else if object.Kind == ast.Var {
			vs, ok := object.Decl.(*ast.ValueSpec)
			if !ok {
				fmt.Println("Error", s, "is ObjectKind Var, but decl is not ValueSpec?")
			}
			l.handleVarDefinition(vs)
		}
	}
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil{
			l.handleFunctionDefinition(fd)
		}
	}
}

func (l linter) handleFunctionDefinition(fd *ast.FuncDecl) {
	name := fd.Name.Name
	if fd.Recv != nil {
		name = fmt.Sprintf("%s.%s", fd.Recv.List[0].Type, name)
	}
	l.handleFunction(name, fd.Pos(), fd.Type.Params.Closing, nil)
}

func (l linter) handleTypeDefinition(typ *ast.TypeSpec) {
	if ft, ok := typ.Type.(*ast.FuncType); ok {
		l.handleFunction(typ.Name.Name, typ.Pos(), ft.Params.Closing, nil)
		return
	}
	if st, ok := typ.Type.(*ast.StructType); ok {
		for _, field := range st.Fields.List {
			if ft, ok := field.Type.(*ast.FuncType); ok {
				startPos := field.Pos()
				l.handleFunction(fmt.Sprintf("%s.%s", typ.Name.Name, field.Names[0].Name), field.Pos(), ft.Params.Closing, &startPos)
			}
		}

	}
}

func (l linter) handleVarDefinition(vari *ast.ValueSpec) {
	for i, value := range vari.Values {
		if fl, ok := value.(*ast.FuncLit); ok {
			l.handleFunction(vari.Names[i].Name, vari.Pos(), fl.Type.Params.Closing, nil)
		}
	}
}

func (l linter) handleFunction(name string, fnPos, paramsClosing token.Pos, startPos *token.Pos) {
	tokenFile := l.fset.File(fnPos)
	line := tokenFile.Line(fnPos)
	length := 0
	if line == tokenFile.LineCount() {
		length = tokenFile.Size() - int(tokenFile.LineStart(line))
	} else {
		length = int((tokenFile.LineStart(line+1) - 1) - tokenFile.LineStart(line))
	}
	if length > 120 {
		fmt.Fprintln(l.out, lengthError(name))
	}
	tabAdjustment := 0
	if startPos != nil {
		tabAdjustment = int(*startPos-tokenFile.LineStart(line))
	}
	lineEnd := tokenFile.Line(paramsClosing)
	if line != lineEnd {
		if int(tokenFile.LineStart(lineEnd))+tabAdjustment != int(paramsClosing) {
			fmt.Fprintln(l.out, formatError(name))
			return
		}
	}
}

func lengthError(name string) error {
	return fmt.Errorf("function declaration too long: %s", name)
}

func formatError(name string) error {
	return fmt.Errorf("params closing not at start of line: %s", name)
}