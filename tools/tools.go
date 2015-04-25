package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonseymour/12coins/lib"
	"os"
)

func main() {
	relabel := false
	reverse := false
	normalize := false
	groupings := false

	flag.BoolVar(&reverse, "reverse", false, "Reverse the solution.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution.")
	flag.BoolVar(&normalize, "normalize", false, "Normalize the solution.")
	flag.BoolVar(&groupings, "groupings", false, "Extract the singletons, pairs and triples.")
	flag.Parse()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var err error
		solver := &lib.Solver{}
		if err = decoder.Decode(&solver); err != nil {
			break
		}

		if reverse {
			if solver, err = solver.Reverse(); err != nil {
				fmt.Fprintf(os.Stderr, "bad solution: %v", err)
			}
		}

		if normalize {
			solver = solver.Normalize()
		}

		if relabel {
			if solver, err = solver.Relabel(); err != nil {
				fmt.Fprintf(os.Stderr, "cannot relabel: %v", err)
			}
		}

		if groupings {
			if solver, err = solver.Groupings(); err != nil {
				fmt.Fprintf(os.Stderr, "bad solution: %v", err)
			}
		}

		encoder.Encode(solver)
	}
}
