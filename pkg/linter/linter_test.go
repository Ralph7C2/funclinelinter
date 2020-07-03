package linter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parse(t *testing.T, src string) (*linter, *ast.File, *bytes.Buffer) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.DeclarationErrors)
	require.NoError(t, err, "parseError")
	buf := bytes.NewBuffer([]byte{})
	return &linter{fset: fset, out: buf}, f, buf
}

func TestFunctionDeclGood(t *testing.T) {
	file := `
package main

func funcWithReallyLongLine(
	someParamWithAReallyLongName thing, more bool, and string, andyet int,
) (int, bool, string, thing) {
	return 0, false, "", thing{}
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Fun, obj.Kind, "ObjectKind is function")
	fd := obj.Decl.(*ast.FuncDecl)
	linter.handleFunctionDefinition(fd)
	assert.Equal(t, "", buf.String(), "output")
}

func TestFunctionDeclLongLine(t *testing.T) {
	file := `
package main

func funcWithReallyLongLine(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, false, "", thing{}
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Fun, obj.Kind, "ObjectKind is function")
	fd := obj.Decl.(*ast.FuncDecl)
	linter.handleFunctionDefinition(fd)
	assert.Equal(t, fmt.Sprintf("%s\n", lengthError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestFunctionDeclWrongFormat(t *testing.T) {
	file := `
package main

func funcWithReallyLongLine(
	someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, false, "", thing{}
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Fun, obj.Kind, "ObjectKind is function")
	fd := obj.Decl.(*ast.FuncDecl)
	linter.handleFunctionDefinition(fd)
	assert.Equal(t, fmt.Sprintf("%s\n", formatError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestFunctionTypeGood(t *testing.T) {
	file := `
package main

type funcWithReallyLongLine = func(
	someParamWithAReallyLongName thing, more bool, and string, andyet int,
) (int, bool, string, thing)
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, "", buf.String(), "output")
}


func TestFunctionTypeLongLine(t *testing.T) {
	file := `
package main

type funcWithReallyLongLine = func(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing)
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, fmt.Sprintf("%s\n", lengthError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestFunctionTypeWrongFormat(t *testing.T) {
	file := `
package main

type funcWithReallyLongLine = func(someParamWithAReallyLongName thing,
	more bool, and string, andyet int,
	) (int, bool, string, thing)
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, fmt.Sprintf("%s\n", formatError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestStructFunctionTypeGood(t *testing.T) {
	file := `
package main

type x struct {
	funcWithReallyLongLine func(
		someParamWithAReallyLongName thing, more bool, and string, andyet int,
	) (int, bool, string, thing)
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("x")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, "", buf.String(), "output")
}

func TestStructFunctionTypeLongLine(t *testing.T) {
	file := `
package main

type x struct {
	funcWithReallyLongLine func(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing)
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("x")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, fmt.Sprintf("%s\n", lengthError("x.funcWithReallyLongLine")), buf.String(), "output")
}

func TestStructFunctionTypeWrongFormat(t *testing.T) {
	file := `
package main

type x struct {
	funcWithReallyLongLine func(
		someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing)
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("x")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Typ, obj.Kind, "ObjectKind is type")
	ts := obj.Decl.(*ast.TypeSpec)
	linter.handleTypeDefinition(ts)
	assert.Equal(t, fmt.Sprintf("%s\n", formatError("x.funcWithReallyLongLine")), buf.String(), "output")
}

func TestVarFunctionTypeGood(t *testing.T) {
	file := `
package main

var funcWithReallyLongLine = func(
	someParamWithAReallyLongName thing, more bool, and string, andyet int,
) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Var, obj.Kind, "ObjectKind is type")
	vs := obj.Decl.(*ast.ValueSpec)
	linter.handleVarDefinition(vs)
	assert.Equal(t, "", buf.String(), "output")
}

func TestVarFunctionTypeLongLine(t *testing.T) {
	file := `
package main

var funcWithReallyLongLine = func(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Var, obj.Kind, "ObjectKind is type")
	vs := obj.Decl.(*ast.ValueSpec)
	linter.handleVarDefinition(vs)
	assert.Equal(t, fmt.Sprintf("%s\n", lengthError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestVarFunctionTypeWrongFormat(t *testing.T) {
	file := `
package main

var funcWithReallyLongLine = func(someParamWithAReallyLongName thing,
	more bool, and string, andyet int,
	) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	obj := f.Scope.Lookup("funcWithReallyLongLine")
	assert.NotNil(t, obj, "function in objects")
	assert.Equal(t, ast.Var, obj.Kind, "ObjectKind is type")
	vs := obj.Decl.(*ast.ValueSpec)
	linter.handleVarDefinition(vs)
	assert.Equal(t, fmt.Sprintf("%s\n", formatError("funcWithReallyLongLine")), buf.String(), "output")
}

func TestMethodTypeGood(t *testing.T) {
	file := `
package main

type x struct {}

func (x) funcWithReallyLongLine(someParamWithAReallyLongName thing,
	more bool, and string, andyet int,
) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil{
			linter.handleFunctionDefinition(fd)
		}
	}
	assert.Equal(t, "", buf.String(), "output")
}

func TestMethodTypeLongLine(t *testing.T) {
	file := `
package main

type x struct {}

func (x) funcWithReallyLongLine(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil{
			linter.handleFunctionDefinition(fd)
		}
	}
	assert.Equal(t, fmt.Sprintf("%s\n", lengthError("x.funcWithReallyLongLine")), buf.String(), "output")
}

func TestMethodTypeWrongFormat(t *testing.T) {
	file := `
package main

type x struct {}

func (x) funcWithReallyLongLine(someParamWithAReallyLongName thing,
	more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, true, "", nil
}
`
	linter, f, buf := parse(t, file)
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil{
			linter.handleFunctionDefinition(fd)
		}
	}
	assert.Equal(t, fmt.Sprintf("%s\n", formatError("x.funcWithReallyLongLine")), buf.String(), "output")
}
func TestLinter_Lint(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	linter := &linter{
		fset: token.NewFileSet(),
		out: buf,
	}
	linter.Lint("./testdata/file.go")
	types := []string{
		"thing.%sFunctionField",
		"thing.%sMethod",
		"%sFunctionType",
		"%sFunctionLiteral",
		"%sFunction",
	}
	var expectedErrors []string
	for _, s := range types {
		expectedErrors = append(expectedErrors, lengthError(fmt.Sprintf(s, "longlyNamed")).Error())
		expectedErrors = append(expectedErrors, formatError(fmt.Sprintf(s, "wronglyFormatted")).Error())
	}
	actualErrors := strings.Split(buf.String(), "\n")
	if actualErrors[len(actualErrors)-1] == "" {
		actualErrors = actualErrors[:len(actualErrors)-1]
	}
	assert.Len(t, actualErrors, len(expectedErrors))
	assert.ElementsMatch(t, expectedErrors, actualErrors)
}