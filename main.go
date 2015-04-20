package main

import (
	"fmt"
	"github.com/jonseymour/12coins/lib"
	"os"
)

// exhaustively test the decision procedure against all possibilities and return those that fail
func main() {

	fail := false
	for i := 0; i < 12; i++ {
		for _, w := range []lib.Weight{lib.Light, lib.Heavy} {
			if err := lib.Test(i, w, decide); err != nil {
				fail = true
				fmt.Fprintf(os.Stderr, "fail: for (%d, %v): %v\n", i, w, err)
			}
		}
	}
	if fail {
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "ok\n")
	}
}
