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
	structure := false
	canonical := false

	flag.BoolVar(&reverse, "reverse", false, "Derive the coins and weights array from the weighings.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution into a indexing form.")
	flag.BoolVar(&normalize, "normalize", false, "Order the coins in each weighing from lowest to highest.")
	flag.BoolVar(&groupings, "groupings", false, "Derive the singletons, pairs and triples from the weighings.")
	flag.BoolVar(&valid, "valid", false, "Only pass valid solutions to stdout.")
	flag.BoolVar(&structure, "structure", false, "Analyse the structure of the weighings.")
	flag.BoolVar(&canonical, "canonical", false, "Permute the weighings into the canonical form.")
	flag.Parse()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var err error
		ok := true
		solution := &lib.Solution{}
		if err = decoder.Decode(&solution); err != nil {
			break
		}

		solution.Decode()

		if reverse {
			if solution, err = solution.Reverse(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: reverse: %v: %v\n", err, solution)
			}
		}

		if normalize {
			solution = solution.Normalize()
		}

		if relabel && ok {
			if solution, err = solution.Relabel(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: relabel: %v: %v\n", err, solution)
			}
		}

		if groupings && ok {
			if solution, err = solution.Groupings(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: groupings: %v: %v\n", err, solution)
			}
		}

		if structure && ok {
			if solution, err = solution.AnalyseStructure(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: structure: %v: %v\n", err, solution)
			}
		}

		if canonical && ok {
			if solution, err = solution.Canonical(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: canonical: %v: %v\n", err, solution)
			}
		}

		if valid {
			if !solution.IsValid() {
				continue
			}
		}

		solution.Encode()
		encoder.Encode(solution)
	}
}
