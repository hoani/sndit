package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hoani/sndit/generate"
)

func main() {
	dir := flag.String("dir", ".", "root directory to scan for sound directories")
	module := flag.String("module", "", "Go module path for generated imports")
	flag.Parse()

	if *module == "" {
		fmt.Fprintln(os.Stderr, "error: -module flag is required")
		os.Exit(1)
	}

	if err := generate.Generate(*dir, *module); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
