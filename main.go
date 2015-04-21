package main

import (
	"fmt"
	"github.com/jonseymour/12coins/lib"
	"os"
)

// exhaustively test the decision procedure against all possibilities and return those that fail
func main() {

	errors := lib.TestAll(decide)
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "%v", e)
		}
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "ok\n")
	}
}
