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
	valid := false

	flag.BoolVar(&reverse, "reverse", false, "Derive the coins and weights array from the weighings.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution into a indexing form.")
	flag.BoolVar(&normalize, "normalize", false, "Order the coins in each weighing from lowest to highest.")
	flag.BoolVar(&groupings, "groupings", false, "Derive the singletons, pairs and triples from the weighings.")
	flag.BoolVar(&valid, "valid", false, "Only pass valid solutions to stdout.")
	flag.Parse()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var err error
		solver := &lib.Solver{}
		if err = decoder.Decode(&solver); err != nil {
			break
		}

		if valid {
			solver.Valid = nil
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

		if valid {
			if solver.Valid == nil {
				if _, err := solver.Reverse(); err != nil {
					continue
				}
			} else if !*solver.Valid {
				continue
			}
		}

		encoder.Encode(solver)
	}
}
