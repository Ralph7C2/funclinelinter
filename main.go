package main

import (
	"flag"
	"fmt"

	"github.com/ralph7c2/funclinelinter/pkg/linter"
)

func main() {
	l := linter.NewLinter()
	flag.Parse()
	for _, arg := range flag.Args() {
		fmt.Println(arg)
		l.Lint(arg)
	}
}
