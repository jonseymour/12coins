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
	reset := false
	invalid := false
	flip := false
	decode := false
	encode := false
	format := false

	flag.BoolVar(&reverse, "reverse", false, "Derive the coins and weights array from the weighings.")
	flag.BoolVar(&flip, "flip", false, "Flip the weighings so that LLL is never a valid weighing.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution into a indexing form.")
	flag.BoolVar(&normalize, "normalize", false, "Order the coins in each weighing from lowest to highest.")
	flag.BoolVar(&groupings, "groupings", false, "Derive the singletons, pairs and triples from the weighings.")
	flag.BoolVar(&valid, "valid", false, "Only pass valid solutions to stdout.")
	flag.BoolVar(&invalid, "invalid", false, "Only pass invalid solutions to stdout.")
	flag.BoolVar(&structure, "structure", false, "Analyse the structure of the weighings.")
	flag.BoolVar(&canonical, "canonical", false, "Permute the weighings into the canonical form.")
	flag.BoolVar(&reset, "reset", false, "Reset the analysis. Implied by reverse, relabel, groupings, structure or canonical.")
	flag.BoolVar(&decode, "decode", false, "Decode a number between 0 and 12!*176 and output the corresponding solution.")
	flag.BoolVar(&encode, "encode", false, "Encode a solution as a number between 0 and 12!*176.")
	flag.BoolVar(&format, "format", false, "Format each solution over multiple lines.")
	flag.Parse()

	if invalid && valid {
		fmt.Fprintf(os.Stderr, "--invalid and --valid are mutually incompatible options - choose one.")
		os.Exit(1)
	}

	if flip {
		reverse = false
	}

	reset = reset || flip || reverse || relabel || groupings || structure || canonical || valid || invalid || encode

	structure = structure || encode

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var err error
		ok := true
		solution := &lib.Solution{}

		if decode {
			var n uint
			if err := decoder.Decode(&n); err != nil {
				break
			}

			if solution, err = lib.DecodeSolution(n); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: decode: %v, %v", err, solution)
			}

		} else {
			if err = decoder.Decode(&solution); err != nil {
				break
			}

			solution.DecodeJSON()

			if reset {
				solution = solution.Reset()
			}
		}

		if reverse {
			if solution, err = solution.Reverse(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: reverse: %v: %v\n", err, solution)
			}
		} else if flip {
			if solution, err = solution.Flip(); err != nil {
				ok = false
				fmt.Fprintf(os.Stderr, "error: flip: %v: %v\n", err, solution)
			}
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

		if normalize {
			solution = solution.Normalize()
		}

		if valid {
			if !solution.IsValid() {
				continue
			}
		}

		if invalid {
			if solution.IsValid() {
				continue
			}
		}

		if encode {
			if ok {
				if n, err := solution.N(); err != nil {
					fmt.Fprintf(os.Stderr, "error: N: %v: %v\n", err, solution)
				} else {
					encoder.Encode(&n)
				}
			}
		} else {
			if format {
				fmt.Fprintf(os.Stdout, "%s", solution.Format())
			} else {
				solution.Encode()
				encoder.Encode(solution)
			}
		}
	}
}
