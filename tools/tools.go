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
	permute := false

	flag.BoolVar(&reverse, "reverse", false, "Reverse the solution.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution.")
	flag.BoolVar(&normalize, "normalize", false, "Normalize the solution.")
	flag.BoolVar(&permute, "permute", false, "Generate permutations of the solution.")
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
				continue
			}
		}

		if normalize {
			solver = solver.Normalize()
		}

		if relabel {
			solver = solver.Relabel()
		}

		if permute {
			ref := solver.Clone()
			for _, m := range []bool{false, true} {
				for _, p := range lib.Permute([]int{0, 1, 2}) {
					clone := ref.Clone()
					clone.Mirror = m
					clone.Permutation = p
					clone.Weighings = [3][2][]int{clone.Weighings[p[0]], clone.Weighings[p[1]], clone.Weighings[p[2]]}
					clone, _ = clone.Reverse()
					if errors := lib.TestAll(clone.Decide); len(errors) != 0 {
						panic(fmt.Errorf("errors: %v", errors))
					}
					encoder.Encode(clone)
				}
			}
			continue
		}

		encoder.Encode(solver)
	}
}
