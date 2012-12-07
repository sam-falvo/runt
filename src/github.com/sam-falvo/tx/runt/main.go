// vim:ts=8:sw=8:noexpandtab:

package main

import (
	"flag"
	"fmt"
	"os"
)

func problem(e error) {
	fmt.Fprintf(os.Stderr, "%s\n", e.Error())
}

func abend() {
	fmt.Fprintf(os.Stderr, "USAGE: %s (options) batch-dir\n", os.Args[0]);
	flag.Usage()
	os.Exit(1)
	panic("os.Exit() didn't os.Exit()!")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "You must at least provide a test batch parameter.\n\n");
		abend()
	}
	s, err := os.Stat(args[0])
	if err != nil {
		problem(err)
		abend()
	}
	if ! s.IsDir() {
		fmt.Fprintf(os.Stderr, "Batch must be a directory.\n\n")
		abend()
	}
}

