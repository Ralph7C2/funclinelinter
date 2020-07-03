package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ralph7c2/funclinelinter/pkg/linter"
)

func main() {
	l := linter.NewLinter()
	flag.Parse()

	fmt.Println(flag.Args())

	for _, arg := range flag.Args() {
		l.Lint(arg)
		out, err := l.Output()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if len(out) != 0 {
			fmt.Println(arg)
			fmt.Println(out)
		}
	}
}
