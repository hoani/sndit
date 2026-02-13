package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hoani/sndit/generate"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	dir := "."
	module := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-dir":
			i++
			if i >= len(args) {
				return errors.New("-dir requires a value")
			}
			dir = args[i]
		case "-module":
			i++
			if i >= len(args) {
				return errors.New("-module requires a value")
			}
			module = args[i]
		}
	}

	if module == "" {
		return errors.New("-module flag is required")
	}

	return generate.Generate(dir, module)
}
