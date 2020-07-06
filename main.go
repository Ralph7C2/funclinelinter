package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ralph7c2/funclinelinter/pkg/linter"
)

func main() {
	fix := flag.Bool("fix", false, "Try to fix lint errors")
	flag.Parse()

	l := linter.NewLinter()
	action := func(string) {}
	if *fix {
		action = linter.Fix
	} else {
		action = l.Lint
	}
	for _, arg := range flag.Args() {
		action(arg)
		if !*fix {
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
}
